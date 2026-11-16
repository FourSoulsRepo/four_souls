package bot

import (
	"context"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// table builds a base-game table: seat 0 is active in its open action
// phase.
func table(t *testing.T, setup engine.SituationSetup) *engine.Game {
	t.Helper()
	g, err := engine.BuildTable(setup, b2.Set)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// best is the standard move of seat p on g with no mistake, the world
// and the bot's random source both from seed.
func best(t *testing.T, g *engine.Game, p engine.PlayerID, seed uint64) engine.Intent {
	t.Helper()
	return bestKnowing(t, g, p, engine.Known{}, seed)
}

// bestKnowing is best, with the seat remembering known.
func bestKnowing(t *testing.T, g *engine.Game, p engine.PlayerID, known engine.Known, seed uint64) engine.Intent {
	t.Helper()
	wf, err := world.For(g, p, known, seed)
	if err != nil {
		t.Fatal(err)
	}
	w, err := wf(0)
	if err != nil {
		t.Fatal(err)
	}
	in, err := standardPick(w, p, w.Allowed(p), rand.New(rand.NewPCG(seed, 9)), 0, false) //nolint:gosec // test moves
	if err != nil {
		t.Fatal(err)
	}
	return in
}

// monster is the object in play of card ref.
func monster(t *testing.T, g *engine.Game, ref engine.CardRef) engine.ObjectID {
	t.Helper()
	for _, m := range g.MonstersInPlay() {
		if g.Object(m).Card == ref {
			return m
		}
	}
	t.Fatalf("no %s in play", ref)
	return 0
}

// seeds are the worlds each built situation is tried on.
var seeds = []uint64{1, 2, 3, 4, 5, 6, 7, 8}

// TestStandardTakesSoul: one soul from a win, a soul monster it can
// likely kill in play: it attacks, and it attacks that monster.
func TestStandardTakesSoul(t *testing.T) {
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players: []engine.SituationPlayer{
				{Character: "isaac", Souls: []engine.CardRef{"monstro", "monstro", "monstro"}},
				{Character: "cain"},
			},
			Monsters: []engine.CardRef{"fly", "conquest"},
		})
		in := best(t, g, 0, seed)
		if in.Kind != engine.IntentAttack {
			t.Fatalf("seed %d: %+v, want an attack", seed, in)
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatal(err)
		}
		for g.Prompt().Kind == engine.PromptPriority { // the attack's priority round
			if _, err := g.Apply(engine.Intent{Player: g.Prompt().Player, Kind: engine.IntentPass}); err != nil {
				t.Fatal(err)
			}
		}
		if g.Prompt().Kind != engine.PromptChoose || g.Choice == nil {
			t.Fatalf("seed %d: no target choice after the attack: %+v", seed, g.Prompt())
		}
		in = best(t, g, 0, seed)
		if in.Kind != engine.IntentChoose || in.Choice >= len(g.Choice.Objects) || g.Choice.Objects[in.Choice] != monster(t, g, "conquest") {
			t.Errorf("seed %d: target %+v, want conquest", seed, in)
		}
	}
}

// TestStandardNoHopelessAttack: at 1 HP, with things to lose, against
// monsters that need two rolls of 6 or five of 4+: a miss is death
// (odds about 3%), so it does not attack; it knows the monster deck's
// top, a monster of 6 HP.
// With full HP and flies in play it does attack.
func TestStandardNoHopelessAttack(t *testing.T) {
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players: []engine.SituationPlayer{{
				Character: "isaac", Damage: 1, Cents: 6,
				Items: []engine.CardRef{"moms_purse", "moms_purse"}, Hand: []engine.CardRef{"soul_heart", "soul_heart"},
			}, {Character: "cain"}},
			Monsters: []engine.CardRef{"little_horn", "mom"},
		})
		// It knows the monster deck's top: no hope there either.
		var known engine.Known
		known.Tops[engine.MonsterDeck] = []engine.CardRef{"satan"}
		if in := bestKnowing(t, g, 0, known, seed); in.Kind == engine.IntentAttack {
			t.Errorf("seed %d: attacked at hopeless odds", seed)
		}
	}
	attacks := 0
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players:  []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}},
			Monsters: []engine.CardRef{"fly", "fly"},
		})
		if best(t, g, 0, seed).Kind == engine.IntentAttack {
			attacks++
		}
	}
	if attacks < len(seeds)/2 {
		t.Errorf("attacked a fly in %d of %d worlds", attacks, len(seeds))
	}
}

// TestStandardKeepsUselessLoot: its only loot card would do nothing now
// (a shield at the end of its turn): it does not play it. A coin card
// it plays.
func TestStandardKeepsUselessLoot(t *testing.T) {
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players: []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"soul_heart"}}, {Character: "cain"}},
		})
		if in := best(t, g, 0, seed); in.Kind == engine.IntentPlayLoot {
			t.Errorf("seed %d: played a loot card for nothing", seed)
		}
		g = table(t, engine.SituationSetup{
			Players: []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"a_nickel"}}, {Character: "cain"}},
		})
		if in := best(t, g, 0, seed); in.Kind != engine.IntentPlayLoot {
			t.Errorf("seed %d: %+v, want the coin card played", seed, in)
		}
	}
}

// TestScore: a win and a loss bound every table; a soul is worth more
// than anything else on the table; the opponents count against.
func TestScore(t *testing.T) {
	g := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain", Souls: []engine.CardRef{"monstro"}}},
	})
	if a, b := Score(g, 0), Score(g, 1); a >= 0 || b <= 0 || a != -b {
		t.Errorf("scores %v and %v: the soul's owner leads, the sum is 0 for two seats", a, b)
	}
	rich := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Cents: 30, Hand: []engine.CardRef{"a_penny", "a_penny", "a_penny"}}, {Character: "cain", Souls: []engine.CardRef{"monstro"}}},
	})
	if Score(rich, 0) >= 0 {
		t.Errorf("30¢ and 3 loot cards outweigh a soul: %v", Score(rich, 0))
	}
	g.Over, g.Winners = true, []engine.PlayerID{1}
	if Score(g, 1) != scoreWin || Score(g, 0) != -scoreWin {
		t.Errorf("game over: %v, %v", Score(g, 1), Score(g, 0))
	}
}

// TestStandardSameSeed: two standard games from equal seeds make equal
// moves, mistakes included.
func TestStandardSameSeed(t *testing.T) {
	play := func() []engine.Intent {
		g := newGame(t, 11, 3)
		var moves []engine.Intent
		for step := range 150 {
			if g.Over {
				break
			}
			p := g.Prompt().Who()[0]
			in, err := Choose(context.Background(), input(t, g, p, Standard, uint64(step))) //nolint:gosec // step is small
			if err != nil {
				t.Fatal(err)
			}
			if _, err := g.Apply(in); err != nil {
				t.Fatalf("step %d: %+v refused: %v", step, in, err)
			}
			moves = append(moves, in)
		}
		return moves
	}
	a, b := play(), play()
	if !slices.EqualFunc(a, b, func(x, y engine.Intent) bool {
		return x.Player == y.Player && x.Kind == y.Kind && x.Choice == y.Choice && slices.Equal(x.Objects, y.Objects)
	}) {
		t.Error("equal seeds gave different moves")
	}
}

// moveStates are positions of a 4-player game with every set, every
// fifth step of a policy game: the moves a standard bot meets.
func moveStates(b *testing.B) []*engine.Game {
	b.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 21, Players: 4, Sets: cards.Sets(), BonusSouls: true})
	if err != nil {
		b.Fatal(err)
	}
	rng := rand.New(rand.NewPCG(21, 1)) //nolint:gosec // test moves
	var out []*engine.Game
	for step := 0; step < 3000 && !g.Over; step++ {
		if step%5 == 0 {
			c, err := g.Clone()
			if err != nil {
				b.Fatal(err)
			}
			out = append(out, c)
		}
		p := g.Prompt().Who()[0]
		in := policyPick(g, p, g.Allowed(p), rng)
		if _, err := g.Apply(in); err != nil {
			b.Fatal(err)
		}
	}
	return out
}

// BenchmarkStandardMove is the time of one standard move (Choose), on
// the positions of moveStates in turn. The copy of world.For is made
// outside the timer, as the room makes it on its own goroutine; the
// imagined world is inside, as the bot asks for it.
func BenchmarkStandardMove(b *testing.B) {
	states := moveStates(b)
	ins := make([]Input, len(states))
	for i, g := range states {
		p := g.Prompt().Who()[0]
		w, err := world.For(g, p, engine.Known{}, uint64(i)) //nolint:gosec // i is small
		if err != nil {
			b.Fatal(err)
		}
		ins[i] = Input{Seat: p, Level: Standard, World: World(w)}
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

// TestJudge: a move must beat doing nothing by more than wAct and the
// noise of its draws; without a way to do nothing, the mean counts.
func TestJudge(t *testing.T) {
	pass := engine.Intent{Kind: engine.IntentPass}
	use := engine.Intent{Kind: engine.IntentActivate}
	cs := []engine.Intent{pass, use}
	nothing := []float64{10, 0, 20, 10}
	for _, c := range []struct {
		name string
		move []float64
		wins bool
	}{
		{"same", []float64{10, 0, 20, 10}, false},
		{"a tenth of a coin", []float64{10.1, 0.1, 20.1, 10.1}, false},
		{"a coin each draw", []float64{11, 1, 21, 11}, true},
		{"a coin on average, by luck", []float64{40, -30, 20, 14}, false},
		{"one draw only", []float64{12}, true},
	} {
		s := judge(cs, [][]float64{nothing, c.move}, noiseOpen, 0)
		if got := s[1] > s[0]; got != c.wins {
			t.Errorf("%s: move %v, nothing %v; move wins %v, want %v", c.name, s[1], s[0], got, c.wins)
		}
	}
	// Two coins on average, one standard error of 0.9: enough on an
	// empty stack, not on a stack.
	gain := [][]float64{nothing, {13, 0, 24, 11}}
	if s := judge(cs, gain, noiseOpen, 0); s[1] <= s[0] {
		t.Errorf("empty stack: %v", s)
	}
	if s := judge(cs, gain, noiseStack, 0); s[1] > s[0] {
		t.Errorf("on a stack: %v", s)
	}
	// On a stack, a coin in every draw but one, where the move wins the
	// game: the win hides the coins' sure gain, unless the noise is
	// clipped.
	won := [][]float64{{10, 0, 20, 10, 10, 0, 20, 10}, {11, 1, 21, 11, 11, 1, 21, scoreWin}}
	if s := judge(cs, won, noiseStack, 0); s[1] > s[0] {
		t.Errorf("a win in one draw, unclipped: %v", s)
	}
	if s := judge(cs, won, noiseStack, noiseClip); s[1] <= s[0] {
		t.Errorf("a win in one draw, clipped: %v", s)
	}
	choose := []engine.Intent{{Kind: engine.IntentChoose, Choice: 0}, {Kind: engine.IntentChoose, Choice: 1}}
	if s := judge(choose, [][]float64{{1, 3}, {2.5}}, noiseStack, 0); s[0] != 2 || s[1] != 2.5 {
		t.Errorf("choices: %v, want the means", s)
	}
}

// TestStandardDigs: no monster in play gives a soul, but a kill lets
// in the next card of the monster deck (Decision 7.9). With good odds
// and an item to lose it attacks; at 1 HP against hard monsters it
// does not.
func TestStandardDigs(t *testing.T) {
	var known engine.Known
	known.Tops[engine.MonsterDeck] = []engine.CardRef{"satan"}
	attacks := 0
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"moms_purse"}}, {Character: "cain"}},
			Monsters: []engine.CardRef{"fly", "black_bony"},
		})
		if bestKnowing(t, g, 0, known, seed).Kind == engine.IntentAttack {
			attacks++
		}
		g = table(t, engine.SituationSetup{
			Players:  []engine.SituationPlayer{{Character: "isaac", Damage: 1, Items: []engine.CardRef{"moms_purse"}}, {Character: "cain"}},
			Monsters: []engine.CardRef{"greedling", "evil_twin"},
		})
		if in := bestKnowing(t, g, 0, known, seed); in.Kind == engine.IntentAttack {
			t.Errorf("seed %d: attacked at bad odds", seed)
		}
	}
	// 7 of 8, not 8: a fly is killed 35 times in 36, but the fixed
	// seed 2 rolls two 1s in a row against it in 2 of its 6 settles, a
	// death each, so that world ends the turn. The seeds are fixed: the
	// count does not move from run to run.
	if attacks < len(seeds)-1 {
		t.Errorf("attacked at good odds in %d of %d worlds", attacks, len(seeds))
	}
}

// TestDigValue: the share of soul cards in the monster deck (soulShare,
// which the policy also gives the deck's top), times wSoul and
// digShare; 0 when the deck holds no soul card.
func TestDigValue(t *testing.T) {
	g := table(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}}})
	deck := g.Decks[engine.MonsterDeck]
	souls := 0
	for _, id := range deck {
		if g.Def(id).Soul > 0 {
			souls++
		}
	}
	if souls == 0 || souls == len(deck) {
		t.Fatalf("%d soul cards of %d: the test needs both kinds", souls, len(deck))
	}
	share := float64(souls) / float64(len(deck))
	if got := soulShare(g); got != share {
		t.Errorf("soul share %v, want %v", got, share)
	}
	if got, want := digValue(g), wSoul*digShare*share; math.Abs(got-want) > 1e-9 {
		t.Errorf("dig %v, want %v", got, want)
	}
	var none []engine.ObjectID
	for _, id := range deck {
		if g.Def(id).Soul == 0 {
			none = append(none, id)
		}
	}
	g.Decks[engine.MonsterDeck] = none
	if got := digValue(g); got != 0 {
		t.Errorf("no soul card left: dig %v, want 0", got)
	}
	g.Decks[engine.MonsterDeck] = nil
	if got := digValue(g); got != 0 {
		t.Errorf("empty deck: dig %v, want 0", got)
	}
}

// TestStandardDiscardsCheap: at the end of its turn the seat holds 11
// loot cards and must discard one (R-TURN-11). By the card values five
// are cheap and six dear: it drops a cheap one, on every world, both by
// standardPick and through Choose, which gives a discard of part of the
// hand to the level (Decision 20). The values are a table of the test's
// own, so the test does not hang on the measured file.
func TestStandardDiscardsCheap(t *testing.T) {
	defer func(v *values) { cardTable = v }(cardTable)
	cardTable = newValues(ValueFile{Cards: []ValueEntry{
		{Key: "a_penny", Kind: KindLoot, Value: 0.1, N: 1000},
		{Key: "a_nickel", Kind: KindLoot, Value: 20, N: 1000},
	}})
	cheap := []engine.CardRef{"a_penny", "a_penny", "a_penny", "a_penny", "a_penny"}
	dear := []engine.CardRef{"a_nickel", "a_nickel", "a_nickel", "a_nickel", "a_nickel", "a_nickel"}
	g := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Hand: append(cheap, dear...)}, {Character: "cain"}},
	})
	if _, err := g.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for range 50 {
		pr := g.Prompt()
		if pr.Kind == engine.PromptDiscard {
			break
		}
		if _, err := g.Apply(engine.Intent{Player: pr.Who()[0], Kind: engine.IntentPass}); err != nil {
			t.Fatal(err)
		}
	}
	if pr := g.Prompt(); pr.Kind != engine.PromptDiscard || pr.Player != 0 || pr.Count != 1 {
		t.Fatalf("prompt %+v, want seat 0 to discard 1", pr)
	}
	check := func(how string, seed uint64, in engine.Intent) {
		t.Helper()
		if in.Kind != engine.IntentDiscard || len(in.Objects) != 1 {
			t.Fatalf("%s seed %d: %+v, want a discard of one card", how, seed, in)
		}
		if c := g.Object(in.Objects[0]).Card; c != "a_penny" {
			t.Errorf("%s seed %d: dropped %s, want a cheap card", how, seed, c)
		}
	}
	for _, seed := range seeds {
		check("standardPick", seed, best(t, g, 0, seed))
		wf, err := world.For(g, 0, engine.Known{}, seed)
		if err != nil {
			t.Fatal(err)
		}
		in, err := Choose(context.Background(), Input{Seat: 0, Level: Standard, Rand: rand.New(rand.NewPCG(seed, 10)), World: World(wf)}) //nolint:gosec // test moves
		if err != nil {
			t.Fatal(err)
		}
		check("Choose", seed, in)
	}
}

// TestStandardPlaysLoot: a card in hand counts wHand of its measured
// value (Decision 20.8.6), so a card that does its usual work is played
// (a_dime, 10¢), and one whose play gives nothing now is kept.
func TestStandardPlaysLoot(t *testing.T) {
	for _, seed := range seeds {
		g := table(t, engine.SituationSetup{
			Players: []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"a_dime"}}, {Character: "cain"}},
		})
		if in := best(t, g, 0, seed); in.Kind != engine.IntentPlayLoot {
			t.Errorf("seed %d: %+v, want the dime played", seed, in)
		}
	}
}

// TestStandardBuysByValue: items are scaled so the average one is worth
// wItem (Decision 20.8.5): with 10¢ and an average item in the shop the
// seat buys it; an item from the bottom of the table it leaves. The
// values are a table of the test's own.
func TestStandardBuysByValue(t *testing.T) {
	defer func(v *values) { cardTable = v }(cardTable)
	cardTable = newValues(ValueFile{Cards: []ValueEntry{
		{Key: "breakfast", Kind: KindPassive, Value: 8, N: 64},
		{Key: "dinner", Kind: KindPassive, Value: 32, N: 64},
		{Key: "lunch", Kind: KindPassive, Value: 8, N: 64},
		{Key: "daddy_haunt", Kind: KindPassive, Value: -16, N: 64},
	}})
	if v := cardTable.cards[kindPassive]["breakfast"]; v != wItem {
		t.Fatalf("the average item is worth %v, want %v", v, wItem)
	}
	for _, c := range []struct {
		item engine.CardRef
		buy  bool
	}{{"breakfast", true}, {"daddy_haunt", false}} {
		for _, seed := range seeds {
			g := table(t, engine.SituationSetup{
				Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, {Character: "cain"}},
				Shop:    []engine.CardRef{c.item},
			})
			in := best(t, g, 0, seed)
			if got := in.Kind == engine.IntentPurchase; got != c.buy {
				t.Errorf("%s, seed %d: %+v, want a purchase %v", c.item, seed, in, c.buy)
			}
		}
	}
}
