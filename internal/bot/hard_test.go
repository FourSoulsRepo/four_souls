package bot

import (
	"context"
	"math/rand/v2"
	"runtime"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// BenchmarkPlayout is the time of one playout (Decision 9.3): a copy of
// a world, a candidate applied, play to the end of the turn, a score;
// on the positions of moveStates in turn. Its inverse is the playouts a
// second of one core.
func BenchmarkPlayout(b *testing.B) {
	states := moveStates(b)
	type pos struct {
		w  *engine.Game
		p  engine.PlayerID
		in engine.Intent
	}
	var ps []pos
	for _, g := range states {
		p := g.Prompt().Who()[0]
		allowed := g.Allowed(p)
		if len(allowed) == 0 {
			continue
		}
		w, err := g.Imagine(p, 1, engine.Known{})
		if err != nil {
			b.Fatal(err)
		}
		ps = append(ps, pos{w, p, complete(w, p, allowed[len(allowed)-1], rand.New(rand.NewPCG(1, 1)))}) //nolint:gosec // test moves
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		x := ps[i%len(ps)]
		if _, err := playout(x.w, x.p, x.in, uint64(i), scoring{lead: true}); err != nil { //nolint:gosec // i is small
			b.Fatal(err)
		}
	}
}

// BenchmarkHardMove is the time, bytes and allocations of one hard move
// of 150 playouts (Choose), on the positions of moveStates in turn that
// have more than one candidate. As in BenchmarkStandardMove, the copy
// of world.For is outside the timer, the imagined worlds inside.
func BenchmarkHardMove(b *testing.B) {
	states := moveStates(b)
	var ins []Input
	for i, g := range states {
		p := g.Prompt().Who()[0]
		if len(g.Allowed(p)) < 2 {
			continue
		}
		w, err := world.For(g, p, engine.Known{}, uint64(i)) //nolint:gosec // i is small
		if err != nil {
			b.Fatal(err)
		}
		ins = append(ins, Input{Seat: p, Level: Hard, World: World(w), Budget: Budget{Playouts: 150}})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		in := ins[i%len(ins)]
		in.Rand = rand.New(rand.NewPCG(uint64(i), 1)) //nolint:gosec // test moves
		if _, err := Choose(context.Background(), in); err != nil {
			b.Fatal(err)
		}
	}
}

// hardBest is the hard move of seat p on g, the worlds and the bot's
// random source from seed, with playouts playouts.
func hardBest(t *testing.T, g *engine.Game, p engine.PlayerID, seed uint64, playouts int) engine.Intent {
	t.Helper()
	wf, err := world.For(g, p, engine.Known{}, seed)
	if err != nil {
		t.Fatal(err)
	}
	in, err := Choose(context.Background(), Input{
		Seat: p, Level: Hard, Rand: rand.New(rand.NewPCG(seed, 9)), World: World(wf), //nolint:gosec // test moves
		Budget: Budget{Playouts: playouts},
	})
	if err != nil {
		t.Fatal(err)
	}
	return in
}

// hardPlayouts is the budget of the built situations.
const hardPlayouts = 150

// passTo passes for every seat before p until p has priority.
func passTo(t *testing.T, g *engine.Game, p engine.PlayerID) {
	t.Helper()
	for range len(g.Players) {
		pr := g.Prompt()
		if pr.Kind == engine.PromptPriority && pr.Player == p {
			return
		}
		if _, err := g.Apply(engine.Intent{Player: pr.Player, Kind: engine.IntentPass}); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatalf("seat %d has no priority: %+v", p, g.Prompt())
}

// answered: in uses the item of card ref (the answer).
func answered(g *engine.Game, in engine.Intent, ref engine.CardRef) bool {
	return in.Kind == engine.IntentActivate && g.Object(in.Objects[0]).Card == ref
}

// attackTable: seat attacker attacks monster with a hit on the stack,
// which kills it; the attacker has 1 HP left, so a miss kills it
// instead. Seat hard holds Godhead, charged: it can turn the roll
// into a 1, a miss. Seats are weak (0 souls, no items) and strong (souls,
// Brimstone for ATK 2 and Mom's Purse).
func attackTable(t *testing.T, seats []engine.SituationPlayer, attacker int, monster engine.CardRef) *engine.Game {
	t.Helper()
	g := table(t, engine.SituationSetup{
		Players:  seats,
		Active:   attacker,
		Monsters: []engine.CardRef{"pin", "horf"},
		Stack:    []engine.SituationStack{{Kind: "attack_roll", DiceRoll: 5, Owner: attacker, Card: monster}},
	})
	return g
}

var (
	weakSeat = engine.SituationPlayer{Character: "isaac", Damage: 1}
	hardSeat = engine.SituationPlayer{Character: "cain", Items: []engine.CardRef{"godhead"}}
)

// strongSeat has souls souls, Brimstone (ATK 2) and Mom's Purse, 1 HP
// left.
func strongSeat(souls int) engine.SituationPlayer {
	s := engine.SituationPlayer{Character: "isaac", Damage: 1, Cents: 8, Items: []engine.CardRef{"brimstone", "moms_purse"}}
	for range souls {
		s.Souls = append(s.Souls, "monstro")
	}
	return s
}

// TestHardAnswersWin: the other seat, one soul from a win, has a hit
// on the stack that kills a soul monster. Hard holds Godhead: it
// answers (Decision 9.7).
func TestHardAnswersWin(t *testing.T) {
	for _, seed := range seeds {
		g := attackTable(t, []engine.SituationPlayer{strongSeat(3), hardSeat}, 0, "pin")
		passTo(t, g, 1)
		if in := hardBest(t, g, 1, seed, hardPlayouts); !answered(g, in, "godhead") {
			t.Errorf("seed %d: %+v, want Godhead on the roll", seed, in)
		}
	}
}

// TestHardSavesAnswer: three seats; hard (seat 2) holds one answer
// that can cancel a monster kill (Decision 7.6). The weak seat kills a
// monster with no soul while the strong seat moves next: hard passes
// and keeps the answer. The strong seat kills a soul monster: hard
// answers, both when the soul wins the game and when it does not.
func TestHardSavesAnswer(t *testing.T) {
	for _, seed := range seeds {
		g := attackTable(t, []engine.SituationPlayer{weakSeat, strongSeat(3), hardSeat}, 0, "horf")
		passTo(t, g, 2)
		if in := hardBest(t, g, 2, seed, hardPlayouts); in.Kind != engine.IntentPass {
			t.Errorf("seed %d: weak seat's kill: %+v, want a pass", seed, in)
		}
		for _, souls := range []int{3, 2} {
			g = attackTable(t, []engine.SituationPlayer{weakSeat, strongSeat(souls), hardSeat}, 1, "pin")
			passTo(t, g, 2)
			if in := hardBest(t, g, 2, seed, hardPlayouts); !answered(g, in, "godhead") {
				t.Errorf("seed %d: strong seat of %d souls kills a soul monster: %+v, want Godhead", seed, souls, in)
			}
		}
	}
}

// TestHardHitsLeader: hard deals 1 damage with Razor Blade; two
// opponents at 1 HP alike but for their souls. It hits the leader.
func TestHardHitsLeader(t *testing.T) {
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players: []engine.SituationPlayer{
				{Character: "cain", Items: []engine.CardRef{"razor_blade"}},
				{Character: "isaac", Damage: 1, Cents: 5, Items: []engine.CardRef{"moms_purse"}, Souls: []engine.CardRef{"monstro", "monstro"}},
				{Character: "isaac", Damage: 1, Cents: 5, Items: []engine.CardRef{"moms_purse"}},
			},
		})
		var use engine.Intent
		for _, in := range g.Allowed(0) {
			if answered(g, in, "razor_blade") {
				use = in
			}
		}
		if use.Kind != engine.IntentActivate {
			t.Fatal("Razor Blade is not usable")
		}
		if _, err := g.Apply(use); err != nil {
			t.Fatal(err)
		}
		if g.Prompt().Kind != engine.PromptChoose || g.Choice == nil || len(g.Choice.Targets) == 0 {
			t.Fatalf("no target choice: %+v", g.Prompt())
		}
		in := hardBest(t, g, 0, seed, hardPlayouts)
		if in.Kind != engine.IntentChoose || in.Choice >= len(g.Choice.Targets) || g.Choice.Targets[in.Choice].Player != 1 {
			t.Errorf("seed %d: %+v, want the leader, seat 1", seed, in)
		}
	}
}

// TestHoldValue: an answer still held is worth a share more than its
// base value with even opponents, and a larger share when one opponent
// is clearly the strongest (Decision 7.6). The answers are what Allowed
// shows on another seat's turn: Godhead and the character's ↷.
func TestHoldValue(t *testing.T) {
	per := func(strong int) float64 {
		g := attackTable(t, []engine.SituationPlayer{weakSeat, strongSeat(strong), hardSeat}, 0, "horf")
		passTo(t, g, 2)
		ans := answers(g, 2, g.Allowed(2))
		godhead := slices.ContainsFunc(ans, func(id engine.ObjectID) bool { return g.Object(id).Card == "godhead" })
		if len(ans) != 2 || !godhead {
			t.Fatalf("answers %v, want Godhead and the character", ans)
		}
		with := scoreAs(g, 2, scoring{lead: true, answers: ans, flat: true})
		without := scoreAs(g, 2, scoring{lead: true, flat: true})
		// With the flat values Godhead's base value is an item's; the
		// character has none. The measured base is the table's.
		return (with - without) / wItemFlat
	}
	even, ahead := per(0), per(3)
	if even < wHold || ahead <= even+wHold {
		t.Errorf("an item answer held is worth %v of its base with even opponents and %v with one 3 souls ahead", even, ahead)
	}
	if got := held(&engine.Game{}, 0, nil, false); got != 0 {
		t.Errorf("no answers: %v held", got)
	}
}

// TestLeadWeights: equal opponents count as their average; one far
// ahead counts nearly alone, and the weights move smoothly with souls.
func TestLeadWeights(t *testing.T) {
	g := table(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "cain"}, {Character: "isaac"}, {Character: "isaac"}}})
	if a, b := Score(g, 0), scoreAs(g, 0, scoring{lead: true}); a != b {
		t.Errorf("equal opponents: average %v, by strength %v", a, b)
	}
	prev := 0.0
	for souls := range 4 {
		var s []engine.CardRef
		for range souls {
			s = append(s, "monstro")
		}
		g = table(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "cain"}, {Character: "isaac", Souls: s}, {Character: "isaac"}}})
		sides := []float64{side(g, 0, 0, false), side(g, 1, 0, false), side(g, 2, 0, false)}
		l := lead(sides, 0)
		if souls > 0 && l <= prev || l < 0 || l > 1 {
			t.Errorf("%d souls ahead: lead %v after %v", souls, l, prev)
		}
		prev = l
		// The leader's part of the opponents: past the average, below
		// the leader alone.
		opp := sides[0] - scoreAs(g, 0, scoring{lead: true})
		if souls > 0 && (opp <= (sides[1]+sides[2])/2 || opp >= sides[1]) {
			t.Errorf("%d souls ahead: opponents count %v; average %v, leader %v", souls, opp, (sides[1]+sides[2])/2, sides[1])
		}
	}
}

// TestHardCancel: a cancelled ctx gives an allowed intent, before
// thinking and while it thinks on a time budget (Decision 9.9).
func TestHardCancel(t *testing.T) {
	g := attackTable(t, []engine.SituationPlayer{weakSeat, strongSeat(2), hardSeat}, 1, "pin")
	passTo(t, g, 2)
	allowed := g.Allowed(2)
	// An hour's budget: only the cancel ends the search, so an answer at
	// all shows the cancel worked; no clock is read (D5 of the review).
	check := func(name string, ctx context.Context) {
		t.Helper()
		in := input(t, g, 2, Hard, 1)
		in.Budget = Budget{Time: time.Hour}
		got, err := Choose(ctx, in)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !slices.ContainsFunc(allowed, func(a engine.Intent) bool { return a.Kind == got.Kind && slices.Equal(a.Objects, got.Objects) }) {
			t.Errorf("%s: %+v is not allowed", name, got)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	check("cancelled before", ctx)
	ctx, cancel = context.WithCancel(context.Background())
	go cancel()
	check("cancelled while thinking", ctx)
}

// TestHardYields: on one core, while hard thinks, another goroutine
// gets the core again and again (Decision 9.9: one core holds many
// tables).
func TestHardYields(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	g := attackTable(t, []engine.SituationPlayer{weakSeat, strongSeat(2), hardSeat}, 1, "pin")
	passTo(t, g, 2)
	// The other goroutine counts each time it gets the one core, and
	// hands it back at once. A search that yields after every playout
	// gives it the core at least once a playout; Go's own preemption,
	// every 10 ms, would give it far fewer. No clock is read.
	var stop atomic.Bool
	var turns atomic.Int64
	done := make(chan struct{})
	go func() {
		defer close(done)
		for !stop.Load() {
			turns.Add(1)
			runtime.Gosched()
		}
	}()
	const playouts = 60
	in := input(t, g, 2, Hard, 1)
	in.Budget = Budget{Playouts: playouts}
	if _, err := Choose(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	got := turns.Load()
	stop.Store(true)
	<-done
	if got < playouts/2 {
		t.Errorf("the other goroutine got the core %d times in a search of %d playouts", got, playouts)
	}
}

// TestHardTimeBudget: a budget of time alone, as the room gives
// (Budget{Time: botThink}; Playouts 0 sets no cap), stops at its time:
// it plays more than one round and answers soon after the time is up.
func TestHardTimeBudget(t *testing.T) {
	g := attackTable(t, []engine.SituationPlayer{weakSeat, strongSeat(2), hardSeat}, 1, "pin")
	passTo(t, g, 2)
	in := input(t, g, 2, Hard, 1)
	const think = 150 * time.Millisecond
	in.Budget = Budget{Time: think}
	start := time.Now()
	if _, err := Choose(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	took := time.Since(start)
	if took < think {
		t.Errorf("answered after %v, before the budget %v ran out", took, think)
	}
	if took > think+raceSlow*100*time.Millisecond {
		t.Errorf("answered after %v; budget %v", took, think)
	}
}
