package server

import (
	"bytes"
	"context"
	"log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	botai "github.com/FourSoulsRepo/four_souls/internal/bot"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// Tests of the bot driver (ADR 030, Decisions 11 – 14; 21.8). They check
// that bot moves are allowed and come in time, never how well bots play.

// botRoom starts a game of the given seats on the test's goroutine: ""
// is a player with a client, else a bot of that level. Nothing runs the
// room: the test hands it its requests (serveBots), so it may look at
// the room between them.
func botRoom(t *testing.T, opts protocol.Options, seats ...string) *Room {
	t.Helper()
	codes := []string{"b2"}
	if err := checkOptions(&opts, codes, len(seats)); err != nil {
		t.Fatal(err)
	}
	rs := make([]roomSeat, len(seats))
	for i, level := range seats {
		if level != "" {
			rs[i] = roomSeat{name: botName(i), bot: level, botN: i}
			continue
		}
		rs[i] = roomSeat{name: "P" + string(rune('0'+i)), token: string(rune('a' + i)), client: &Client{conn: newMemConn()}}
	}
	r, err := NewRoom(engineSetup(7, len(seats), codes, opts), opts, rs)
	if err != nil {
		t.Fatal(err)
	}
	r.second, r.botThink = time.Millisecond, 2*time.Millisecond
	r.recInfo.game = "g1"
	t.Cleanup(r.stopThinkers)
	r.play(characters(setsOf(codes))[:len(seats)])
	wake(t, r)
	return r
}

// wake starts the thinkers the room noted, as Run does after each
// request.
func wake(t *testing.T, r *Room) {
	t.Helper()
	r.startThinkers(t.Context())
}

// serveBots runs the room until done says so: it hands the room each
// bot's move, and a player asked to answer passes or gives the forced
// answer at once.
func serveBots(t *testing.T, r *Room, done func() bool) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for wake(t, r); !done(); wake(t, r) {
		if time.Now().After(deadline) {
			t.Fatal("the room did not get there in time")
		}
		if r.playerAnswers() {
			continue
		}
		select {
		case req := <-r.in:
			r.handle(req)
		case <-time.After(5 * time.Second):
			t.Fatal("no bot move came")
		}
	}
}

// playerAnswers gives the forced answer for the first waiting player
// seat (a pass or the end of the turn); false: no player is asked.
func (r *Room) playerAnswers() bool {
	if r.game == nil || r.game.Over || r.paused() || r.rbVote != nil {
		return false
	}
	for _, p := range r.game.Prompt().Who() {
		if r.seats[p].bot != "" || r.seats[p].client == nil {
			continue
		}
		in, ok := r.forcedIntent(p)
		if !ok {
			continue
		}
		events, err := r.apply(in)
		if err != nil {
			continue
		}
		r.after(events, true)
		return true
	}
	return false
}

// thinking: the bot seat has an open thinker.
func thinking(r *Room, seat int) func() bool {
	return func() bool { return len(r.thinkers) > seat && r.thinkers[seat] != nil }
}

// noBotMove hands the room what comes in for a while: an answer a bot
// sent before its thinker was dropped is dropped, so no step is made.
func noBotMove(t *testing.T, r *Room, when string) {
	t.Helper()
	step := r.step
	end := time.After(30 * time.Millisecond)
	for {
		select {
		case req := <-r.in:
			r.handle(req)
			if r.step != step {
				t.Fatalf("a bot moved %s", when)
			}
		case <-end:
			return
		}
	}
}

// TestBotsNeverAway: a bot seat is never away; a game of a player and
// bots is not paused (ADR 030, Decision 11.2).
func TestBotsNeverAway(t *testing.T) {
	r := botRoom(t, protocol.Options{}, "", protocol.BotStandard, protocol.BotHard)
	if r.paused() || len(r.away()) != 0 {
		t.Fatalf("away %v: a bot seat is never away", r.away())
	}
	seats := r.seatViews()
	if seats[1].Bot != protocol.BotStandard || seats[2].Bot != protocol.BotHard || seats[2].BotN != 2 {
		t.Errorf("seats %+v", seats)
	}
}

// TestBotsWaitWhilePausedVotingOver: no bot move while the game is
// paused, during a rollback vote, or after the game end; an answer that
// comes then is dropped, and the bot thinks again once the game goes on
// (Decision 12.1, 12.9).
func TestBotsWaitWhilePausedVotingOver(t *testing.T) {
	r := botRoom(t, protocol.Options{Rollback: protocol.RollbackSoft}, "", protocol.BotStandard)
	serveBots(t, r, thinking(r, 1))

	// Paused: the player left.
	th := r.thinkers[1]
	step := r.step
	ann := r.seats[0].client
	ann.seat = 0
	r.drop(ann)
	if !r.paused() || r.thinkers[1] != nil || !closed(th.done) {
		t.Fatal("the pause did not drop the thinker")
	}
	r.handle(request{bot: &botAnswer{seat: 1, id: th.id, step: th.step, intent: engine.Intent{Kind: engine.IntentPass}}})
	if r.step != step {
		t.Fatal("a bot moved while paused")
	}
	noBotMove(t, r, "while paused")
	back := &Client{conn: newMemConn()}
	r.attach(back, "a", "")
	wake(t, r)
	if r.paused() || r.thinkers[1] == nil {
		t.Fatal("the bot does not think again after the pause")
	}

	// A rollback vote is open.
	th, step = r.thinkers[1], r.step
	r.rbVote = &softVote{author: 1}
	r.broadcast(nil)
	if r.thinkers[1] != nil || !closed(th.done) {
		t.Fatal("the vote did not drop the thinker")
	}
	r.handle(request{bot: &botAnswer{seat: 1, id: th.id, step: th.step, intent: engine.Intent{Kind: engine.IntentPass}}})
	if r.step != step {
		t.Fatal("a bot moved during a rollback vote")
	}
	noBotMove(t, r, "during a rollback vote")
	r.endVote(protocol.VoteCancelled)
	wake(t, r)
	if r.thinkers[1] == nil {
		t.Fatal("the bot does not think again after the vote")
	}

	// The game is over.
	th, step = r.thinkers[1], r.step
	r.game.Over = true
	r.broadcast(nil)
	if r.thinkers[1] != nil || !closed(th.done) {
		t.Fatal("the game end did not drop the thinker")
	}
	r.handle(request{bot: &botAnswer{seat: 1, id: th.id, step: th.step, intent: engine.Intent{Kind: engine.IntentPass}}})
	if r.step != step {
		t.Fatal("a bot moved after the game end")
	}
	noBotMove(t, r, "after the game end")
}

// TestBotsRollbackThinksAgain: a rollback of a bot's move drops the
// bot's open thinker; its late answer is dropped, and the bot thinks
// again on the restored game (Decision 12.8, 12.9).
func TestBotsRollbackThinksAgain(t *testing.T) {
	r := botRoom(t, protocol.Options{Rollback: protocol.RollbackGod}, "", protocol.BotHard)
	botStep := func() int {
		return slices.IndexFunc(r.live, func(s liveStep) bool { return s.rec.Intent.Player == 1 && r.isTarget(s) })
	}
	serveBots(t, r, func() bool { return botStep() >= 0 && thinking(r, 1)() })
	old := r.thinkers[1]
	if err := r.rollBackTo(0, botStep()); err != nil {
		t.Fatal(err)
	}
	wake(t, r)
	if !closed(old.done) {
		t.Fatal("the rollback did not drop the open thinker")
	}
	again := r.thinkers[1]
	if again == nil || again.id == old.id || again.step != r.step {
		t.Fatalf("thinker after the rollback %+v; before %+v", again, old)
	}
	step, n := r.step, len(r.live)
	r.handle(request{bot: &botAnswer{seat: 1, id: old.id, step: old.step, intent: engine.Intent{Kind: engine.IntentPass}}})
	if r.step != step {
		t.Fatal("the dropped thinker's answer was applied")
	}
	serveBots(t, r, func() bool { return len(r.live) > n })
	if r.live[n].rec.Intent.Player != 1 {
		t.Fatalf("the first step after the rollback %+v; want the bot's new move", r.live[n].rec.Intent)
	}
}

// TestBotsVoteHumansOnly: a kick vote and a Soft vote count only the
// players; bots never vote (ADR 030, Decision 11.3).
func TestBotsVoteHumansOnly(t *testing.T) {
	r := botRoom(t, protocol.Options{Rollback: protocol.RollbackSoft}, "", "", protocol.BotStandard, protocol.BotHard)
	if got := r.voters(2); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("Soft voters on a bot's move %v, want the players only", got)
	}
	if got := r.voters(0); !slices.Equal(got, []int{1}) {
		t.Errorf("Soft voters on a player's move %v, want the other player only", got)
	}
	ann, bob := r.seats[0].client, r.seats[1].client
	ann.seat, bob.seat = 0, 1
	r.drop(bob)
	if !slices.Equal(r.away(), []int{1}) {
		t.Fatalf("away %v", r.away())
	}
	msg, err := protocol.Encode(protocol.TypeVote, 1, protocol.Vote{Kick: true})
	if err != nil {
		t.Fatal(err)
	}
	env, err := protocol.Decode(msg)
	if err != nil {
		t.Fatal(err)
	}
	r.vote(ann, env)
	if !r.seats[1].kicked || r.paused() {
		t.Fatal("one vote of the one connected player kicks; bots do not count")
	}
}

// TestBotsRefusedIntent: an intent the engine refuses gives the forced
// answer, and the server logs it; so does a bot with no move (Decision
// 12.10).
func TestBotsRefusedIntent(t *testing.T) {
	var buf syncBuffer
	prev := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(prev)
	r := botRoom(t, protocol.Options{Rollback: protocol.RollbackGod}, "", protocol.BotStandard)
	for _, a := range []botAnswer{
		{intent: engine.Intent{Kind: engine.IntentAttack, Objects: []engine.ObjectID{99999}}},
		{err: context.Canceled},
	} {
		serveBots(t, r, thinking(r, 1))
		th := r.thinkers[1]
		forced, ok := r.forcedIntent(1)
		if !ok {
			t.Fatal("no forced answer")
		}
		forced.Player = 1
		n := len(r.live)
		a.seat, a.id, a.step = 1, th.id, th.step
		r.handle(request{bot: &a})
		if len(r.live) <= n || !intentEqual(r.live[n].rec.Intent, forced) {
			t.Fatalf("live %d → %d; want the forced answer %+v", n, len(r.live), forced)
		}
	}
	out := buf.String()
	if !strings.Contains(out, "Bot 1: the engine refused") || !strings.Contains(out, "Bot 1 has no move") {
		t.Errorf("log %q: the refusal and the missing move are not logged", out)
	}
}

// closed reports whether a thinker's context ended.
func closed(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// intentEqual compares intents by their JSON fields that matter here.
func intentEqual(a, b engine.Intent) bool {
	return a.Player == b.Player && a.Kind == b.Kind && slices.Equal(a.Objects, b.Objects) && a.Choice == b.Choice
}

// syncBuffer is a log buffer safe for goroutines.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p) //nolint:wrapcheck // a test buffer
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// TestBotsPace: a bot's pass or answer comes 0.7 of the room's second
// after the prompt; an action in its own turn 1 to 2 (Decision 13).
func TestBotsPace(t *testing.T) {
	r := botRoom(t, protocol.Options{}, "", protocol.BotStandard)
	r.second = 40 * time.Millisecond
	r.stopThinkers()
	r.drive() // think again at the new pace
	wake(t, r)
	for range 8 {
		serveBots(t, r, thinking(r, 1))
		th := r.thinkers[1]
		began := time.Now()
		var req request
		for req.bot == nil || req.bot.id != th.id {
			select {
			case req = <-r.in:
			case <-time.After(5 * time.Second):
				t.Fatal("no bot move")
			}
		}
		took := time.Since(began)
		own := r.game.Turn.Active == 1 && r.game.PromptFor(1).Kind == engine.PromptPriority && req.bot.intent.Kind != engine.IntentPass
		least := r.botSecond(botAnswerPace) * 3 / 4 // the thinker started a little before began
		if own {
			least = r.botSecond(botActionPace) * 3 / 4
		}
		if took < least || took > r.botSecond(botActionPace+1)+time.Second {
			t.Errorf("a bot move after %v; own turn action %v", took, own)
		}
		r.handle(req)
	}
}

// TestBotsBanAndPick: bots ban on their ban turns and pick from their
// draft offers; the game starts once the player has done both
// (Decision 14).
func TestBotsBanAndPick(t *testing.T) {
	h := startHubWith(t, time.Millisecond)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 3, Options: &protocol.Options{BanRounds: 1, Picking: protocol.PickDraft, DraftSize: 2}})
	for i := 1; i < 3; i++ {
		ann.seatBot(i, protocol.BotHard)
	}
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	var last protocol.Setup
	for {
		msg := <-ann.conn.out
		e, err := protocol.Decode(msg)
		if err != nil {
			t.Fatal(err)
		}
		if e.Type == protocol.TypeUpdate {
			break
		}
		if e.Type != protocol.TypeSetup {
			continue
		}
		var s protocol.Setup
		if err := e.Unpack(&s); err != nil {
			t.Fatal(err)
		}
		last = s
		switch {
		case s.Phase == protocol.PhaseBan && s.Turn == 0:
			ann.say(protocol.TypeBan, protocol.BanCard{Card: s.Pool[0]})
		case s.Phase == protocol.PhasePick && len(s.Offers) > 0 && !s.Picked[0]:
			for i := 1; i < 3; i++ {
				if !s.Picked[i] {
					t.Errorf("bot %d has not picked yet", i)
				}
			}
			ann.say(protocol.TypePick, protocol.PickCard{Card: s.Offers[0]})
		}
	}
	var seats []int
	for _, b := range last.Banned {
		seats = append(seats, b.Seat)
	}
	slices.Sort(seats)
	if !slices.Equal(seats, []int{0, 1, 2}) {
		t.Errorf("bans by seat %v, want one each", seats)
	}
}

// TestBotsGame plays seeded games of one player and three bots, of both
// levels, to their end (21.8, Done when). The player only passes and
// ends turns, or gives the forced answer. The bot driver, with its
// copies for the thinkers, never holds the room over the stall limit.
func TestBotsGame(t *testing.T) {
	seeds := make([]uint64, 20)
	for i := range seeds {
		seeds[i] = uint64(1000 + i) //nolint:gosec // a small game number
	}
	if raceEnabled || testing.Short() {
		seeds = []uint64{1012, 1013} // the race detector slows whole games ten times
	}
	games := len(seeds)
	// A few games at a time: twenty at once load the machine so much
	// that any request may stall, a player's as well as a bot's.
	slots := make(chan struct{}, 4)
	var mu sync.Mutex
	var most time.Duration
	for i, seed := range seeds {
		levels := []string{protocol.BotStandard, protocol.BotHard, protocol.BotStandard}
		if i%2 == 1 {
			levels = []string{protocol.BotHard, protocol.BotStandard, protocol.BotHard}
		}
		t.Run(strconv.FormatUint(seed, 10), func(t *testing.T) {
			t.Parallel()
			slots <- struct{}{}
			defer func() { <-slots }()
			work, turns := playBotsGame(t, seed, levels)
			mu.Lock()
			most = max(most, work)
			mu.Unlock()
			t.Logf("seed %d, bots %v: %d turns; the bot driver held the room at most %v", seed, levels, turns, work)
			if work > stallLimit {
				t.Errorf("the bot driver held the room %v, over the stall limit %v", work, stallLimit)
			}
		})
	}
	t.Cleanup(func() { t.Logf("the bot driver held the room at most %v in %d games", most, games) })
}

// playBotsGame plays one game of a player and bots of the given levels
// on a fresh hub. It returns the longest time the bot driver held the
// room, and the game's turns.
func playBotsGame(t *testing.T, seed uint64, levels []string) (time.Duration, int) {
	t.Helper()
	h := NewHub()
	h.second, h.seed, h.botThink = time.Millisecond, seed, 2*time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go h.Run(ctx)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: len(levels) + 1})
	ann.next(protocol.TypeTable, nil)
	for i, level := range levels {
		ann.seatBot(i+1, level)
	}
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	const limit = 2 * time.Minute
	start := time.Now()
	for {
		var u protocol.Update
		ann.next(protocol.TypeUpdate, &u)
		if u.Pause != nil {
			t.Fatalf("seed %d: the game paused: %+v", seed, u.Pause)
		}
		if u.View.Over {
			return time.Duration(ann.c.room.Load().botWork.Load()), u.View.Turn.Number
		}
		if time.Since(start) > limit {
			t.Fatalf("seed %d: no end after %v, turn %d", seed, limit, u.View.Turn.Number)
		}
		if in, ok := passOrForced(u); ok {
			ann.say(protocol.TypeIntent, protocol.Intent{Intent: in})
		}
	}
}

// passOrForced is the player's answer to an update: a pass, the end of
// the turn, or the first allowed intent, a discard of the fewest cards.
func passOrForced(u protocol.Update) (engine.Intent, bool) {
	if len(u.Allowed) == 0 {
		return engine.Intent{}, false
	}
	for _, kind := range []engine.IntentKind{engine.IntentPass, engine.IntentEndTurn} {
		if i := slices.IndexFunc(u.Allowed, func(x engine.Intent) bool { return x.Kind == kind }); i >= 0 {
			return u.Allowed[i], true
		}
	}
	in := u.Allowed[0]
	if in.Kind == engine.IntentDiscard && u.View.Waiting.Count <= len(in.Objects) {
		in.Objects = in.Objects[:u.View.Waiting.Count]
	}
	return in, true
}

// TestBotsMustAttackFallback: with the response timer off, a bot that
// must attack (no pass, no end of turn: R-ABIL-12, R-TURN-09) and has
// no move, or gives one the engine refuses, still moves: the forced
// answer has none, so an allowed attack goes, and the game goes on.
func TestBotsMustAttackFallback(t *testing.T) {
	for _, a := range []botAnswer{
		{err: context.Canceled},
		{intent: engine.Intent{Kind: engine.IntentAttack, Objects: []engine.ObjectID{99999}}},
	} {
		r := botRoom(t, protocol.Options{ResponseTimer: 0, Rollback: protocol.RollbackGod}, "", protocol.BotStandard) // God mode keeps the steps to look at
		serveBots(t, r, func() bool {
			g := r.game
			return thinking(r, 1)() && g.Turn.Active == 1 && slices.ContainsFunc(g.Allowed(1), func(x engine.Intent) bool { return x.Kind == engine.IntentEndTurn }) &&
				slices.ContainsFunc(g.Allowed(1), func(x engine.Intent) bool { return x.Kind == engine.IntentAttack })
		})
		r.game.Turn.MustAttacks = 1
		if _, ok := r.forcedIntent(1); ok {
			t.Fatalf("a forced answer while the bot must attack: allowed %+v", r.game.Allowed(1))
		}
		th, step := r.thinkers[1], r.step
		a.seat, a.id, a.step = 1, th.id, th.step
		r.handle(request{bot: &a})
		if r.step == step || r.live[len(r.live)-1].rec.Intent.Player != 1 {
			t.Fatalf("no bot step after %+v: the bot stands", a)
		}
		if k := r.live[len(r.live)-1].rec.Intent.Kind; k != engine.IntentAttack {
			t.Errorf("the fallback move %v, want an attack", k)
		}
	}
}

// TestBotsPanic: a bot panic, in the room's copy (world.For) or in the
// thinker (an imagined world), is logged with its stack; the bot's
// fallback move goes and the game goes on (Decision 12.12).
func TestBotsPanic(t *testing.T) {
	for _, tc := range []struct {
		name     string
		worldFor func(*engine.Game, engine.PlayerID, engine.Known, uint64) (world.Func, error)
	}{
		{"copy", func(*engine.Game, engine.PlayerID, engine.Known, uint64) (world.Func, error) {
			panic("a copy that breaks")
		}},
		{"world", func(*engine.Game, engine.PlayerID, engine.Known, uint64) (world.Func, error) {
			return func(int) (*engine.Game, error) { panic("a world that breaks") }, nil
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf syncBuffer
			prev := log.Writer()
			log.SetOutput(&buf)
			defer log.SetOutput(prev)
			r := botRoom(t, protocol.Options{Rollback: protocol.RollbackGod}, "", protocol.BotStandard) // God mode keeps the steps to look at
			r.stopThinkers()
			r.worldFor = tc.worldFor
			r.drive()
			botSteps := func() int {
				n := 0
				for _, l := range r.live {
					if l.rec.Intent.Player == 1 {
						n++
					}
				}
				return n
			}
			before := botSteps()
			serveBots(t, r, func() bool { return botSteps() >= before+3 })
			out := buf.String()
			for _, want := range []string{"bot panic: seed 7", ", seat 1: a ", "goroutine ", "Bot 1 has no move: server: bot panic"} {
				if !strings.Contains(out, want) {
					t.Errorf("log %q: no %q", out, want)
				}
			}
		})
	}
}

// TestBotLevel: a seat's protocol level gives the bot package's level;
// the server knows standard and hard only (ADR 030, Decision 1.2).
func TestBotLevel(t *testing.T) {
	for level, want := range map[string]botai.Level{protocol.BotStandard: botai.Standard, protocol.BotHard: botai.Hard} {
		if got, err := botLevel(level); err != nil || got != want {
			t.Errorf("botLevel(%q) = %q, %v; want %q", level, got, err, want)
		}
	}
	for _, level := range []string{string(botai.Random), string(botai.Policy), string(botai.Flat), ""} {
		if _, err := botLevel(level); err == nil {
			t.Errorf("botLevel(%q) gives a level", level)
		}
	}
}
