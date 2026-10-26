package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestLemonMishap(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"lemon_mishap"}}, seat("cain")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb.Activate(0, "lemon_mishap", 0)
	tb.G.ForceRolls(1, 3, 3, 3) // a 1 misses: each monster takes 1; 3s hit (no trigger)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fatty", "ask.rest_in_slot_order")
	monstro, _ := tb.G.Monsters[0].TopOf()
	if d := tb.G.Object(monstro).Damage; d != 1 {
		t.Errorf("Monstro damage = %d, want 1 from the roll of 1", d)
	}
	// It ends with the turn.
	tb.EndTurn()
	if len(tb.G.Lingering) != 0 {
		t.Error("the lingering ability outlived the turn")
	}
}

// Placebo uses Lemon Mishap's ↷: the lingering ability is Lemon Mishap's
// (a simulation found this).
func TestLemonMishapThroughPlacebo(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"placebo"}}, seat("cain", "lemon_mishap")},
		Monsters: []engine.CardRef{"monstro", "fatty"},
	})
	tb.Activate(0, "placebo", 0)
	tb.G.ForceRolls(1, 3, 3, 3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fatty", "ask.rest_in_slot_order")
	monstro, _ := tb.G.Monsters[0].TopOf()
	if d := tb.G.Object(monstro).Damage; d != 1 {
		t.Errorf("Monstro damage = %d, want 1", d)
	}
}
