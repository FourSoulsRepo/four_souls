package server

import (
	"context"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// setupTable seats n players at a game with the given options and marks
// them ready.
func setupTable(t *testing.T, h *Hub, n int, opts protocol.Options) []*player {
	t.Helper()
	host, _ := arrive(t, h, "P0", "")
	host.say(protocol.TypeCreate, protocol.Create{Seats: n, Options: &opts})
	var tb protocol.Table
	host.next(protocol.TypeTable, &tb)
	players := []*player{host}
	for i := 1; i < n; i++ {
		p, _ := arrive(t, h, "P"+strconv.Itoa(i), "")
		p.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
		p.next(protocol.TypeTable, nil)
		players = append(players, p)
	}
	for _, p := range players {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	return players
}

// dealt returns the characters of the game's first update.
func dealt(p *player) []engine.CardRef {
	p.t.Helper()
	var u protocol.Update
	p.next(protocol.TypeUpdate, &u)
	var out []engine.CardRef
	for _, pl := range u.View.Players {
		out = append(out, pl.Character.Card)
	}
	return out
}

var baseCharacters = []engine.CardRef{"blue_baby", "cain", "eden", "eve", "isaac", "judas", "lazarus", "lilith", "maggy", "samson", "the_forgotten"}

// GS-02: host bans are never dealt.
func TestHostBans(t *testing.T) {
	h := startHub(t)
	// Ban all but Isaac and Samson: those two must be dealt.
	bans := slices.DeleteFunc(slices.Clone(baseCharacters), func(c engine.CardRef) bool { return c == "isaac" || c == "samson" })
	players := setupTable(t, h, 2, protocol.Options{HostBans: bans})
	got := dealt(players[0])
	for _, c := range got {
		if slices.Contains(bans, c) {
			t.Errorf("banned %s was dealt", c)
		}
	}
	if len(got) != 2 || !slices.Contains(got, "isaac") || !slices.Contains(got, "samson") {
		t.Fatalf("dealt %v, want Isaac and Samson", got)
	}
}

// GS-13: the host may ban any card; banned cards are in no deck, and
// the record keeps the list.
func TestHostBansAnyCard(t *testing.T) {
	dir := t.TempDir()
	h := NewHub()
	h.records = dir
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	bans := []engine.CardRef{"delirium", "stoney", "belly_button", "isaac"}
	players := setupTable(t, h, 2, protocol.Options{HostBans: bans})
	var u protocol.Update
	players[0].next(protocol.TypeUpdate, &u)
	for _, pl := range u.View.Players {
		if slices.Contains(bans, pl.Character.Card) {
			t.Errorf("banned %s was dealt", pl.Character.Card)
		}
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
	if !slices.Equal(s.Banned, []string{"delirium", "stoney", "belly_button", "isaac"}) {
		t.Fatalf("record bans %v", s.Banned)
	}
	setup := engine.Setup{Seed: s.Seed, Players: s.Players, BonusSouls: s.BonusSouls, Banned: bans}
	for _, code := range s.Sets {
		set, _ := cards.Find(code)
		setup.Sets = append(setup.Sets, set)
	}
	for _, c := range s.Characters {
		setup.Characters = append(setup.Characters, engine.CardRef(c))
	}
	g, _, err := engine.NewGame(setup)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range g.Objects {
		if slices.Contains(bans, o.Card) {
			t.Errorf("banned %s is in the game", o.Card)
		}
	}
}

// GS-04: ban rounds go in snake order; bans leave the pool.
func TestBanRounds(t *testing.T) {
	h := startHub(t)
	players := setupTable(t, h, 3, protocol.Options{BanRounds: 2})
	want := []int{0, 1, 2, 2, 1, 0}
	var banned []engine.CardRef
	for i, seat := range want {
		var s protocol.Setup
		players[0].next(protocol.TypeSetup, &s)
		if s.Phase != protocol.PhaseBan || s.Turn != seat || s.Round != i/3+1 {
			t.Fatalf("ban %d: phase %s turn %d round %d, want turn %d", i, s.Phase, s.Turn, s.Round, seat)
		}
		if i == 0 {
			players[1].say(protocol.TypeBan, protocol.BanCard{Card: s.Pool[0]})
			if code := players[1].failure(); code != protocol.ErrNotYourTurn {
				t.Errorf("out of turn: %s", code)
			}
			players[0].say(protocol.TypeBan, protocol.BanCard{Card: "fly"})
			if code := players[0].failure(); code != protocol.ErrBadCard {
				t.Errorf("not a character: %s", code)
			}
		}
		players[seat].say(protocol.TypeBan, protocol.BanCard{Card: s.Pool[0]})
		banned = append(banned, s.Pool[0])
	}
	for _, c := range dealt(players[0]) {
		if slices.Contains(banned, c) {
			t.Errorf("banned %s was dealt", c)
		}
	}
}

// GS-04: when the ban timer runs out, that player makes no ban.
func TestBanTimer(t *testing.T) {
	h := startHubWith(t, time.Millisecond)
	players := setupTable(t, h, 2, protocol.Options{BanRounds: 1, BanTimer: 5})
	var s protocol.Setup
	players[0].next(protocol.TypeSetup, &s)
	if s.Deadline == 0 {
		t.Error("no deadline shown")
	}
	// Nobody bans: both turns time out and the characters are dealt.
	if got := dealt(players[0]); len(got) != 2 {
		t.Fatalf("dealt %v", got)
	}
}

// GS-01, GS-03: a draft offers each player their own characters to pick
// from.
func TestDraft(t *testing.T) {
	h := startHub(t)
	players := setupTable(t, h, 2, protocol.Options{Picking: protocol.PickDraft, DraftSize: 3})
	var offers [2][]engine.CardRef
	for i, p := range players {
		var s protocol.Setup
		p.next(protocol.TypeSetup, &s)
		if s.Phase != protocol.PhasePick || len(s.Offers) != 3 {
			t.Fatalf("seat %d: phase %s offers %v", i, s.Phase, s.Offers)
		}
		offers[i] = s.Offers
	}
	for _, c := range offers[0] {
		if slices.Contains(offers[1], c) {
			t.Fatalf("offers overlap: %v and %v", offers[0], offers[1])
		}
	}
	players[0].say(protocol.TypePick, protocol.PickCard{Card: offers[1][0]})
	if code := players[0].failure(); code != protocol.ErrBadCard {
		t.Errorf("picking another's offer: %s", code)
	}
	players[0].say(protocol.TypePick, protocol.PickCard{Card: offers[0][2]})
	players[1].say(protocol.TypePick, protocol.PickCard{Card: offers[1][1]})
	got := dealt(players[1])
	if got[0] != offers[0][2] || got[1] != offers[1][1] {
		t.Errorf("characters %v, want the picks %s and %s", got, offers[0][2], offers[1][1])
	}
}

func TestOptionChecks(t *testing.T) {
	h := startHub(t)
	ann, _ := arrive(t, h, "Ann", "")
	for name, o := range map[string]protocol.Options{
		"draft of 9":       {Picking: protocol.PickDraft, DraftSize: 9},
		"unknown picking":  {Picking: "auction"},
		"4 ban rounds":     {BanRounds: 4},
		"1 second timer":   {BanRounds: 1, BanTimer: 1},
		"not in the sets":  {HostBans: []engine.CardRef{"nope"}},
		"too few to draft": {Picking: protocol.PickDraft, DraftSize: 5, BanRounds: 3},
	} {
		ann.say(protocol.TypeCreate, protocol.Create{Seats: 4, Options: &o})
		if code := ann.failure(); code != protocol.ErrBadSetup {
			t.Errorf("%s: %s", name, code)
		}
	}
}

// GS-07: deck ratios are on by default when the sets have the cards for
// them, off with too few cards or when the host turns them off; the
// record keeps whether a game used them.
func TestDeckRatiosOption(t *testing.T) {
	dir := t.TempDir()
	h := NewHub()
	h.records = dir
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	all := []string{"b2", "g2", "fsp2", "r"}
	hello := func(name string) *player {
		p := &player{t: t, conn: newMemConn()}
		p.c = h.Connect(p.conn)
		p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: name, Role: protocol.RolePlayer, Sets: all})
		p.next(protocol.TypeWelcome, nil)
		return p
	}
	for name, c := range map[string]struct {
		sets []string
		opts *protocol.Options
		off  bool
	}{
		"every set":        {sets: all},
		"Base Game only":   {sets: []string{"b2"}, off: true},
		"turned off":       {sets: all, opts: &protocol.Options{NoDeckRatios: true}, off: true},
		"detailed, on":     {sets: all, opts: &protocol.Options{}},
		"too few for 2":    {sets: []string{"b2", "r"}, opts: &protocol.Options{}, off: true},
		"bans make it few": {sets: all, opts: &protocol.Options{HostBans: []engine.CardRef{"lost_soul"}}, off: true},
	} {
		p := hello("Ann")
		p.say(protocol.TypeCreate, protocol.Create{Seats: 2, Sets: c.sets, Options: c.opts})
		var tb protocol.Table
		p.next(protocol.TypeTable, &tb)
		if tb.Options == nil || tb.Options.NoDeckRatios != c.off {
			t.Errorf("%s: options %+v, want no_deck_ratios %v", name, tb.Options, c.off)
		}
	}

	host := hello("Bo")
	host.say(protocol.TypeCreate, protocol.Create{Seats: 2, Sets: all})
	var tb protocol.Table
	host.next(protocol.TypeTable, &tb)
	guest := hello("Cy")
	guest.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	guest.next(protocol.TypeTable, nil)
	for _, p := range []*player{host, guest} {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	host.next(protocol.TypeUpdate, nil)
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
	if !rec.Header.Setup.DeckRatios {
		t.Error("the record header lacks deck_ratios")
	}
}

// ADR 023: a back face is never dealt, so the pool holds only fronts.
func TestBackFacesAreNotDealt(t *testing.T) {
	pool := characters(cards.Sets())
	if !slices.Contains(pool, "the_enigma") {
		t.Fatal("The Enigma is not in the pool")
	}
	for _, set := range cards.Sets() {
		for _, d := range set.Cards {
			if d.Back != "" && slices.Contains(pool, d.Back) {
				t.Errorf("%s: back %s is in the pool", d.Ref, d.Back)
			}
		}
	}
}
