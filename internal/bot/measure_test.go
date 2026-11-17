package bot

import (
	"math"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// fixedWorld is a World that is always a copy of g: the measure's
// positions are built tables, nothing hidden is dealt.
func fixedWorld(g *engine.Game) World {
	return func(int) (*engine.Game, error) { return g.Clone() }
}

// TestPlayGainKnown: a card that gives 5¢ and nothing else is worth,
// played, what the flat score says 5¢ more and the card out of the hand
// are worth: the measure adds nothing of its own (Decision 20.4).
func TestPlayGainKnown(t *testing.T) {
	before := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"a_nickel"}}, {Character: "cain"}},
	})
	after := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Cents: 5}, {Character: "cain"}},
	})
	sc := scoring{flat: true}
	want := scoreAs(after, 0, sc) - scoreAs(before, 0, sc)
	var nickel engine.ObjectID = -1
	for _, id := range before.Players[0].Hand {
		nickel = id
	}
	gain, ok, err := PlayGain(fixedWorld(before), 0, nickel, []uint64{1, 2, 3, 4})
	if err != nil || !ok {
		t.Fatalf("PlayGain: ok %v, %v", ok, err)
	}
	if math.Abs(gain-want) > 0.25*wCent {
		t.Errorf("a nickel gains %.3f, want %.3f (5¢ less its place in the hand)", gain, want)
	}
}

// TestItemPlayoutRepeats: the same world on the same seed scores the
// same, so two worlds that differ by nothing differ by 0: a card's
// value is the difference its presence makes and nothing else.
func TestItemPlayoutRepeats(t *testing.T) {
	g := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}},
	})
	for seed := uint64(1); seed <= 4; seed++ {
		a, err := ItemPlayout(fixedWorld(g), 0, seed)
		if err != nil {
			t.Fatal(err)
		}
		b, err := ItemPlayout(fixedWorld(g), 0, seed)
		if err != nil {
			t.Fatal(err)
		}
		if a != b {
			t.Errorf("seed %d: %v and %v", seed, a, b)
		}
	}
}

// TestValuesShrink: a card's value is weighed against its kind's mean
// by precision (Decision 20.8): v' = (v/se² + m/τ²) / (1/se² + 1/τ²),
// τ² the spread of the kind's values less their mean se²; then items
// are scaled so their mean is wItem.
func TestValuesShrink(t *testing.T) {
	v := newValues(ValueFile{Cards: []ValueEntry{
		{Key: "breakfast", Kind: KindPassive, Value: 10, SE: 1, N: 50},
		{Key: "lunch", Kind: KindPassive, Value: 30, SE: 1, N: 50},
	}})
	m := 20.0
	tau2 := (math.Pow(10-m, 2)+math.Pow(30-m, 2))/1 - 1 // sample variance (n-1 = 1) less the mean se²
	shrunk := (10/1.0 + m/tau2) / (1/1.0 + 1/tau2)
	scale := wItem * 2 / (10 + 30)
	if got, want := v.cards[kindPassive]["breakfast"], shrunk*scale; math.Abs(got-want) > 1e-9 {
		t.Errorf("breakfast %v, want %v", got, want)
	}
}
