package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSiblingRivalry(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain", "sibling_rivalry")},
		Monsters: []engine.CardRef{"pooter"},
	})
	tb.G.ForceRolls(6, 6) // Isaac hits, Cain hits: the pooter (2 HP) dies
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Activate(1, "sibling_rivalry", 0, foe, "pooter") // in answer to the declaration (S-YUG)
	if top(tb, 0) == "pooter" {
		t.Fatal("the pooter should be dead")
	}
	if len(tb.G.Players[0].Hand) != 1 || len(tb.G.Players[1].Hand) != 1 {
		t.Errorf("both gain its reward (loot 1): hands %d and %d", len(tb.G.Players[0].Hand), len(tb.G.Players[1].Hand))
	}
}
