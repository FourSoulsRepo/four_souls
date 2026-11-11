package challenge_test

import (
	"math/rand/v2"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Trick or Treat (R-TRICK, 14.32). Every attack on a monster starts
// with the rule card's roll(s): with Trick up, a 6 only makes the
// monster's rewards 1¢, so tests that look elsewhere roll 6 first.

const (
	hauntNorm  = "mischievous_the_haunt_norm"
	hauntUltra = "mischievous_the_haunt_ultra"
	trickCard  = "trick_treat_normhardultra"
	treatCard  = "treat_normhardultra"
)

// bareSet has a monster without rewards, which the real sets lack
// (Q-54.5).
var bareSet = engine.CardSet{Code: "trick_trial", Name: "trick_trial", Cards: []engine.CardDef{
	{Ref: "bare_mob", Kind: engine.MonsterCard, HP: 1, DC: 2, ATK: 1, Outside: true},
}}

func trickSets() []engine.CardSet { return append(sets(), bareSet) }

// trick builds a Trick or Treat table at difficulty d: Isaac (active)
// and Cain unless the seats name others; the monsters given, else
// leech and fly.
func trick(t *testing.T, d string, isaac, cain engine.SituationPlayer, monsters ...engine.CardRef) *enginetest.Table {
	t.Helper()
	if len(monsters) == 0 {
		monsters = []engine.CardRef{"leech", "fly"}
	}
	return trickOn(t, d, engine.SituationSetup{Players: []engine.SituationPlayer{isaac, cain}, Monsters: monsters})
}

func trickOn(t *testing.T, d string, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	for i, c := range []engine.CardRef{"isaac", "cain"} {
		if s.Players[i].Character == "" {
			s.Players[i].Character = c
		}
	}
	tb := enginetest.NewSetup(t, s, trickSets()...)
	if err := tb.G.StartChallenge(&challenge.TrickOrTreat, d); err != nil {
		t.Fatal(err)
	}
	// Plain monsters on the deck: a refill draws no curse or event.
	for range 6 {
		tb.G.PutOnDeck(engine.MonsterDeck, "bare_mob")
	}
	return tb
}

func tricks(tb *enginetest.Table) int {
	return tb.G.Object(boss(tb)).CountersOf(challenge.TrickCounter)
}

// setTricks leaves n trick counters on the boss.
func setTricks(tb *enginetest.Table, n int) {
	tb.G.AddNamedCounters(boss(tb), challenge.TrickCounter, n-tricks(tb))
}

func fear(tb *enginetest.Table, p engine.PlayerID) int {
	return tb.G.PlayerCounters(p, challenge.FearCounter)
}

// flipToTreat turns the rule card to Treat, as a yes would.
func flipToTreat(tb *enginetest.Table) {
	tb.T.Helper()
	if !tb.G.Flip(tb.G.Challenge.Card) || tb.G.CardOf(tb.G.Challenge.Card) != treatCard {
		tb.T.Fatal("the rule card did not flip to Treat")
	}
}

// settle throws, passes and answers until the stack is empty: a
// trigger order takes its first option; any other question takes the
// next choice, or stays open when none is left.
func settle(tb *enginetest.Table, choices ...string) []engine.Event {
	tb.T.Helper()
	var out []engine.Event
	g := tb.G
	for range 1000 {
		w := g.Prompt()
		var in engine.Intent
		switch {
		case g.Over:
			return out
		case w.Kind == engine.PromptChoose && g.Choice != nil && g.Choice.Purpose == engine.ChooseTriggerOrder:
			in = engine.Intent{Player: w.Player, Kind: engine.IntentChoose, Choice: 0}
		case w.Kind == engine.PromptChoose && len(choices) > 0:
			i := g.OptionIndex(choices[0])
			if i < 0 {
				tb.T.Fatalf("%q is not an option of %v", choices[0], tb.Options())
			}
			choices = choices[1:]
			in = engine.Intent{Player: w.Player, Kind: engine.IntentChoose, Choice: i}
		case w.Kind == engine.PromptThrow:
			in = engine.Intent{Player: w.Player, Kind: engine.IntentThrow}
		case w.Kind == engine.PromptPriority && (len(g.Stack) > 0 || g.Attack.On || g.Purchase.On):
			in = engine.Intent{Player: w.Player, Kind: engine.IntentPass}
		default:
			if len(choices) > 0 {
				tb.T.Fatalf("choices %v left over", choices)
			}
			return out
		}
		ev, err := g.Apply(in)
		if err != nil {
			tb.T.Fatal(err)
		}
		out = append(out, ev...)
	}
	tb.T.Fatal("never settled")
	return nil
}

// asking reports whether a question waits for an answer.
func asking(tb *enginetest.Table) bool {
	return !tb.G.Over && tb.G.Prompt().Kind == engine.PromptChoose
}

// attack has the active player attack monster with the given rolls,
// then settles; no question may stay open.
func attack(tb *enginetest.Table, monster engine.CardRef, rolls ...int) []engine.Event {
	tb.T.Helper()
	ev := attackAsk(tb, monster, rolls...)
	if asking(tb) {
		tb.T.Fatalf("a question is open: %v %+v", tb.Options(), *tb.G.Choice)
	}
	return ev
}

// attackAsk is attack, but a question may stay open.
func attackAsk(tb *enginetest.Table, monster engine.CardRef, rolls ...int) []engine.Event {
	tb.T.Helper()
	tb.G.ForceRolls(rolls...)
	tb.ThrowPending()
	if _, err := tb.G.Apply(engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentAttack}); err != nil {
		tb.T.Fatal(err)
	}
	return settle(tb, string(monster))
}

// hp gives player p n more HP this turn.
func hp(tb *enginetest.Table, p engine.PlayerID, n int) {
	tb.G.Boosts = append(tb.G.Boosts, engine.Boost{Stat: engine.StatPlayerHP, Player: p, Amount: n})
}

// damageTo is the damage instances player p took in events.
func damageTo(ev []engine.Event, p engine.PlayerID) []int {
	var out []int
	for _, e := range ev {
		if e.Kind == engine.EvDamaged && e.Player == p && e.Object == 0 {
			out = append(out, e.Amount)
		}
	}
	return out
}

// rollsOf counts the rolls the card's abilities put on the stack.
func rollsOf(tb *enginetest.Table, ev []engine.Event, card engine.CardRef) int {
	n := 0
	for _, e := range ev {
		if e.Kind == engine.EvStackAdded && e.Label.Key.String() == "event.roll" && tb.G.CardOf(e.Object) == card {
			n++
		}
	}
	return n
}

// R-TRICK-02 to R-TRICK-04, ADR 028: per difficulty the boss sits in a
// final boss slot with its trick counters; one rule card, Trick side
// up; the view shows both. R-TRICK-16, Q-49: the regular The Haunt
// stays in the game.
func TestTrickSetup(t *testing.T) {
	if c, ok := cards.FindChallenge("trick_or_treat"); !ok || c != &challenge.TrickOrTreat {
		t.Fatal("no trick_or_treat challenge")
	}
	for d, want := range map[string]struct {
		boss   engine.CardRef
		tricks int
	}{
		engine.DifficultyNormal: {hauntNorm, 12},
		engine.DifficultyHard:   {"mischievous_the_haunt_hard", 14},
		engine.DifficultyUltra:  {hauntUltra, 14},
	} {
		s := engine.Setup{Seed: 3, Players: 2, Sets: append(cards.Sets(), challenge.Set), Challenge: &challenge.TrickOrTreat, Difficulty: d}
		g, _, err := engine.NewGame(s)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := g.Bosses[0].TopOf()
		if !g.IsFinalBoss(b) || g.Object(b).Card != want.boss || g.Object(b).CountersOf(challenge.TrickCounter) != want.tricks {
			t.Errorf("%s: boss %s with %d trick counters", d, g.Object(b).Card, g.Object(b).CountersOf(challenge.TrickCounter))
		}
		if g.HP(b) != 4 || g.Evasion(b) != 5 || g.Def(b).ATK != 1 || g.Def(b).Soul != 1 || g.Rewards(b)[0].Kind != engine.RewardTreasure {
			t.Errorf("%s: boss stats HP %d DC %d", d, g.HP(b), g.Evasion(b))
		}
		card := g.Challenge.Card
		if g.Object(card).Flipped || g.CardOf(card) != trickCard {
			t.Errorf("%s: rule card %s flipped %v", d, g.CardOf(card), g.Object(card).Flipped)
		}
		v := g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0})
		if v.Challenge.Card.Card != trickCard || v.Challenge.Card.Flipped || v.Challenge.ID != "trick_or_treat" {
			t.Errorf("%s: view card %+v", d, v.Challenge.Card)
		}
		if c := v.Bosses[0].Top.Counters; len(c) != 1 || c[0] != (engine.Counter{Name: "trick", Count: want.tricks}) {
			t.Errorf("%s: view counters %+v", d, c)
		}
		haunt := false
		for _, o := range g.Objects {
			haunt = haunt || (o.Card == "the_haunt" && o.Zone.Kind != engine.ZoneGone)
		}
		if !haunt {
			t.Errorf("%s: the regular The Haunt is out of the game", d)
		}
		check(t, g)
	}
}

// R-TRICK-01, R-COOP-09: seeded whole games end. Players with souls and
// few trick counters who always attack, the boss first, win; random
// players lose by the timer. Replays and saves keep the challenge.
func TestTrickWholeGames(t *testing.T) {
	for _, c := range []struct {
		seed   uint64
		d      string
		greedy bool
		won    bool
	}{
		{7, engine.DifficultyNormal, true, true},
		{20, engine.DifficultyHard, true, true},
		{7, engine.DifficultyUltra, true, true},
		{1, engine.DifficultyNormal, false, false},
		{2, engine.DifficultyHard, false, false},
		{3, engine.DifficultyUltra, false, false},
	} {
		g := playTrick(t, c.seed, c.d, c.greedy)
		if !g.Over || g.Challenge.Lost == c.won {
			t.Errorf("seed %d %s: over %v, lost %v (%s)", c.seed, c.d, g.Over, g.Challenge.Lost, g.Challenge.Reason.Key)
		}
	}
}

// playTrick plays a whole game. Greedy players start with 3 monster
// souls between them and 4 trick counters on the boss; they attack, the
// boss first, and never flip. Random players start as dealt.
func playTrick(t *testing.T, seed uint64, d string, greedy bool) *engine.Game {
	t.Helper()
	s := engine.Setup{Seed: seed, Players: 2, Sets: append(cards.Sets(), challenge.Set), Challenge: &challenge.TrickOrTreat, Difficulty: d}
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	if greedy {
		for i := 0; i < 3; {
			g.MillTop(engine.MonsterDeck)
			pile := g.Discards[engine.MonsterDeck]
			if top := pile[len(pile)-1]; g.Kind(top) == engine.MonsterCard {
				g.SoulFromDiscard(engine.MonsterDeck, g.Object(top).Card, engine.PlayerID(i%2))
				i++
			}
		}
		b, _ := g.Bosses[0].TopOf()
		g.AddNamedCounters(b, challenge.TrickCounter, 4-g.Object(b).CountersOf(challenge.TrickCounter))
	}
	base, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	rng := rand.New(rand.NewPCG(seed, 7)) //nolint:gosec // test moves
	var steps []engine.StepRecord
	for step := 0; step < 20000 && !g.Over; step++ {
		p := g.Prompt()
		in := choose(g, g.Allowed(p.Player), rng, greedy)
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, serr := g.Step(in)
		if serr != nil {
			t.Fatalf("step %d: %+v: %v", step, in, serr)
		}
		steps = append(steps, st)
		check(t, g)
	}
	r, err := engine.ReplayFrom(base, steps, s.Sets...)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	a, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Save()
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Error("the replay ends elsewhere")
	}
	if _, err := engine.Load(a, s.Sets...); err != nil {
		t.Fatalf("load: %v", err)
	}
	return g
}
