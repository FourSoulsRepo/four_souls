package server

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// krampus is a Krampus challenge at the given difficulty (ADR 028).
func krampus(d string) *protocol.Options {
	return &protocol.Options{Challenge: "krampus", Difficulty: d}
}

// checkKrampus: the view is a Krampus challenge at d, with its final
// boss and no bonus souls (R-COOP-02, R-CHAL-01).
func checkKrampus(t *testing.T, v engine.GameView, d string) {
	t.Helper()
	c := v.Challenge
	if c == nil || c.ID != "krampus" || c.Difficulty != d || len(c.Presents) != 4 {
		t.Fatalf("challenge %+v, want krampus %s with 4 presents", c, d)
	}
	if len(v.Bosses) != 1 || v.Bosses[0].Top == nil || v.Bosses[0].Top.Card != engine.CardRef("pilfering_krampus_"+map[string]string{"normal": "norm", "hard": "hard", "ultra": "ultra"}[d]) {
		t.Errorf("bosses %+v", v.Bosses)
	}
	if len(v.BonusSouls) != 0 {
		t.Errorf("bonus souls %+v in a challenge", v.BonusSouls)
	}
}

// ADR 028, Decision 10: a Krampus table has 2 seats; bonus souls go off
// whatever the option said; the game list names the challenge; the
// game starts as a challenge with the challenge's own set.
func TestChallengeCreate(t *testing.T) {
	h := startHub(t)
	host, _ := arrive(t, h, "Ann", "")
	opts := krampus(engine.DifficultyHard)
	host.say(protocol.TypeCreate, protocol.Create{Seats: 2, Options: opts})
	var tb protocol.Table
	host.next(protocol.TypeTable, &tb)
	if tb.Options == nil || tb.Options.Challenge != "krampus" || tb.Options.Difficulty != engine.DifficultyHard || !tb.Options.NoBonusSouls {
		t.Fatalf("table options %+v", tb.Options)
	}
	if len(tb.Sets) != 1 || tb.Sets[0] != "b2" {
		t.Errorf("table sets %v: the challenge set is not listed", tb.Sets)
	}
	guest, _ := arrive(t, h, "Bo", "")
	var games protocol.Games
	guest.next(protocol.TypeGames, &games)
	if len(games.Games) != 1 || games.Games[0].Challenge != "krampus" || games.Games[0].Difficulty != engine.DifficultyHard {
		t.Fatalf("game list %+v", games.Games)
	}
	guest.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	guest.next(protocol.TypeTable, nil)
	for _, p := range []*player{host, guest} {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	var u protocol.Update
	host.next(protocol.TypeUpdate, &u)
	checkKrampus(t, u.View, engine.DifficultyHard)

	// The started game keeps its challenge in the list.
	cy, _ := arrive(t, h, "Cy", "")
	cy.next(protocol.TypeGames, &games)
	if len(games.Games) != 1 || !games.Games[0].Started || games.Games[0].Challenge != "krampus" || games.Games[0].Difficulty != engine.DifficultyHard {
		t.Errorf("started game in the list %+v", games.Games)
	}

	// Bonus souls asked for are turned off.
	dee, _ := arrive(t, h, "Dee", "")
	dee.say(protocol.TypeCreate, protocol.Create{Seats: 2, Options: &protocol.Options{Challenge: "krampus", Difficulty: engine.DifficultyNormal, NoBonusSouls: false}})
	dee.next(protocol.TypeTable, &tb)
	if !tb.Options.NoBonusSouls {
		t.Error("bonus souls stay on in a challenge")
	}
}

// ADR 028, Decision 10.2: refusals carry a protocol error and a message
// key.
func TestChallengeRefusals(t *testing.T) {
	h := startHub(t)
	ann, _ := arrive(t, h, "Ann", "")
	for name, c := range map[string]struct {
		seats int
		opts  protocol.Options
		key   string
	}{
		"3 seats":                     {3, *krampus(engine.DifficultyNormal), "error.challenge_seats"},
		"4 seats":                     {4, *krampus(engine.DifficultyNormal), "error.challenge_seats"},
		"unknown challenge":           {2, protocol.Options{Challenge: "nope", Difficulty: engine.DifficultyNormal}, "error.no_challenge"},
		"no difficulty":               {2, protocol.Options{Challenge: "krampus"}, "error.no_difficulty"},
		"bad difficulty":              {2, *krampus("easy"), "error.no_difficulty"},
		"competitive is no level":     {2, *krampus("comp"), "error.no_difficulty"},
		"difficulty and no challenge": {2, protocol.Options{Difficulty: engine.DifficultyHard}, "error.no_difficulty"},
		"ban the boss":                {2, protocol.Options{Challenge: "krampus", Difficulty: engine.DifficultyNormal, HostBans: []engine.CardRef{"pilfering_krampus_norm"}}, "error.server"},
		"ban a rule card":             {2, protocol.Options{Challenge: "krampus", Difficulty: engine.DifficultyNormal, HostBans: []engine.CardRef{"how_the_krampus_stole_christmas_hard"}}, "error.server"},
	} {
		ann.say(protocol.TypeCreate, protocol.Create{Seats: c.seats, Options: &c.opts})
		var e protocol.Error
		ann.next(protocol.TypeError, &e)
		if e.Code != protocol.ErrBadSetup || e.Text == nil || e.Text.Key.String() != c.key {
			t.Errorf("%s: %s %+v, want %s %s", name, e.Code, e.Text, protocol.ErrBadSetup, c.key)
		}
	}
}

// ADR 028: the other options keep working in a challenge: every set,
// the deck ratios, the room deck, a draft, ban rounds, rollback, random
// seats.
func TestChallengeWithOptions(t *testing.T) {
	h := startHub(t)
	all := []string{"b2", "g2", "fsp2", "r"}
	hello := func(name string) *player {
		p := &player{t: t, conn: newMemConn()}
		p.c = h.Connect(p.conn)
		p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: name, Role: protocol.RolePlayer, Sets: all})
		p.next(protocol.TypeWelcome, nil)
		return p
	}
	host := hello("Ann")
	opts := protocol.Options{
		Challenge: "krampus", Difficulty: engine.DifficultyUltra, Picking: protocol.PickDraft, BanRounds: 1,
		RoomDeck: true, Rollback: protocol.RollbackSoft, RandomSeats: true, HostBans: []engine.CardRef{"isaac"},
	}
	host.say(protocol.TypeCreate, protocol.Create{Seats: 2, Sets: all, Options: &opts})
	var tb protocol.Table
	host.next(protocol.TypeTable, &tb)
	o := tb.Options
	if o == nil || o.NoDeckRatios || !o.RoomDeck || !o.NoBonusSouls || o.Picking != protocol.PickDraft {
		t.Fatalf("options %+v: want deck ratios, the room deck, no bonus souls, a draft", o)
	}
	guest := hello("Bo")
	guest.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	guest.next(protocol.TypeTable, nil)
	players := []*player{host, guest}
	for _, p := range players {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	// One ban each, in turn, then each picks the first offer.
	for range 2 {
		var su protocol.Setup
		for su.Phase != protocol.PhaseBan {
			host.next(protocol.TypeSetup, &su)
		}
		who := host
		if su.Seats[su.Turn].Name != "Ann" {
			who = guest
		}
		who.say(protocol.TypeBan, protocol.BanCard{Card: su.Pool[0]})
	}
	for _, p := range players {
		var su protocol.Setup
		for su.Phase != protocol.PhasePick || len(su.Offers) == 0 {
			p.next(protocol.TypeSetup, &su)
		}
		p.say(protocol.TypePick, protocol.PickCard{Card: su.Offers[0]})
	}
	var u protocol.Update
	host.next(protocol.TypeUpdate, &u)
	checkKrampus(t, u.View, engine.DifficultyUltra)
	if len(u.View.Rooms) == 0 {
		t.Error("no room slot in a challenge with the room deck")
	}
	if u.RollbackMode != protocol.RollbackSoft {
		t.Errorf("rollback mode %q", u.RollbackMode)
	}
}

// ADR 028, Decision 8.5: a lost challenge ends the game: the record
// ends finished with no winners, and the room is over (RP-05).
func TestChallengeLossFinishes(t *testing.T) {
	dir := t.TempDir()
	opts := *krampus(engine.DifficultyNormal)
	if msg := checkChallenge(&opts, 2); !msg.IsZero() {
		t.Fatal(msg)
	}
	if err := checkOptions(&opts, []string{"b2"}, 2); err != nil {
		t.Fatal(err)
	}
	setup := engineSetup(3, 2, []string{"b2"}, opts)
	room, err := NewRoom(setup, opts, []roomSeat{{name: "Ann", token: "a"}, {name: "Bo", token: "b"}})
	if err != nil {
		t.Fatal(err)
	}
	room.recDir, room.path = dir, filepath.Join(dir, "g1"+record.Ext)
	room.recInfo = recordInfo{game: "g1", sets: []string{"b2"}}
	room.play([]engine.CardRef{"isaac", "cain"})
	if room.game.Challenge == nil || room.over.Load() {
		t.Fatal("no challenge game")
	}
	room.game.LoseChallenge(engine.ReasonTimer.Msg())
	room.after(nil, false)
	if !room.over.Load() {
		t.Fatal("the lost game is not over")
	}
	rec, err := record.ReadFile(room.path)
	if err != nil {
		t.Fatal(err)
	}
	if rec.End == nil || !rec.End.Finished || len(rec.End.Winners) != 0 {
		t.Errorf("record end %+v, want finished with no winners", rec.End)
	}
	s := rec.Header.Setup
	if s.Challenge != "krampus" || s.Difficulty != engine.DifficultyNormal || s.BonusSouls || len(s.Sets) != 1 {
		t.Errorf("record setup %+v", s)
	}
}

// ADR 007, ADR 028 Decision 10.5: a challenge record's header names the
// challenge and its difficulty; the game replays from it with every
// checksum.
func TestChallengeRecordReplays(t *testing.T) {
	dir := t.TempDir()
	h := NewHub()
	h.records = dir
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	players := startTableWith(t, h, protocol.Create{Seats: 2, Options: krampus(engine.DifficultyNormal)})
	for range 40 { // whoever is asked answers with an attack, else the first allowed intent
		for _, p := range players {
			up, ok := p.latest()
			if !ok || len(up.Allowed) == 0 || up.View.Over {
				continue
			}
			in := up.Allowed[0]
			for _, a := range up.Allowed {
				if a.Kind == engine.IntentAttack {
					in = a
				}
			}
			if in.Kind == engine.IntentDiscard {
				in.Objects = in.Objects[:up.View.Waiting.Count]
			}
			p.say(protocol.TypeIntent, protocol.Intent{Intent: in})
		}
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	<-stopped

	files, err := filepath.Glob(filepath.Join(dir, "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	rec, err := record.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	s := rec.Header.Setup
	if s.Challenge != "krampus" || s.Difficulty != engine.DifficultyNormal || s.BonusSouls {
		t.Fatalf("record setup %+v", s)
	}
	if len(rec.Steps) == 0 {
		t.Fatal("the record has no steps")
	}
	steps := make([]engine.StepRecord, len(rec.Steps))
	for i, st := range rec.Steps {
		if jerr := json.Unmarshal(st.Intent, &steps[i].Intent); jerr != nil {
			t.Fatal(jerr)
		}
		steps[i].Checksum = st.Checksum
	}
	g, err := engine.Replay(replaySetup(t, s), steps)
	if err != nil {
		t.Fatalf("the challenge record does not replay: %v", err)
	}
	if g.Challenge == nil || g.Challenge.ID != "krampus" {
		t.Errorf("replayed challenge %+v", g.Challenge)
	}
}

// N-11: a saved challenge loads again and goes on as the same challenge.
func TestChallengeSaveAndLoad(t *testing.T) {
	h := startHubWithSaves(t, t.TempDir())
	players := startTableWith(t, h, protocol.Create{Seats: 2, Options: krampus(engine.DifficultyHard)})
	var before protocol.Update
	players[0].next(protocol.TypeUpdate, &before)
	checkKrampus(t, before.View, engine.DifficultyHard)
	players[0].say(protocol.TypeSave, struct{}{})
	var saved protocol.Saved
	players[0].next(protocol.TypeSaved, &saved)
	players[1].next(protocol.TypeSaved, nil)

	players[0].say(protocol.TypeLoad, protocol.Load(saved))
	var tb protocol.Table
	players[0].next(protocol.TypeTable, &tb)
	if tb.Options == nil || tb.Options.Challenge != "krampus" || tb.Options.Difficulty != engine.DifficultyHard {
		t.Fatalf("loaded table options %+v", tb.Options)
	}
	players[1].say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	players[1].next(protocol.TypeTable, nil)
	for _, p := range players {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	var after protocol.Update
	players[0].next(protocol.TypeUpdate, &after)
	checkKrampus(t, after.View, engine.DifficultyHard)
	if after.View.Challenge.Timer != before.View.Challenge.Timer || after.View.Bosses[0].Top.Card != before.View.Bosses[0].Top.Card {
		t.Errorf("the loaded challenge differs: %+v, was %+v", after.View.Challenge, before.View.Challenge)
	}
}

// ADR 019, ADR 028 Decision 10.6: a challenge kept over a server
// restart goes on as the same challenge.
func TestChallengeRestart(t *testing.T) {
	dir := t.TempDir()
	bots, _ := stoppedGameWith(t, dir, krampus(engine.DifficultyUltra))
	next := startKeeping(t, dir)
	reconnect(next, bots...)
	for _, b := range bots {
		u := b.untilBack()
		checkKrampus(t, u.View, engine.DifficultyUltra)
	}
}
