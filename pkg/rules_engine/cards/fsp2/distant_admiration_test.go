package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestDistantAdmiration(t *testing.T) {
	plays := 0
	tb := setup(t, engine.SituationSetup{
		Players:   []engine.SituationPlayer{{Character: "isaac", Items: items("distant_admiration")}, {Character: "cain", Hand: items("three_cents", "bomb")}},
		LootPlays: &plays,
	})
	tb.Activate(0, "distant_admiration", 0, foe, "three_cents")
	if c := tb.G.Players[0].Cents; c != 3 {
		t.Errorf("cents = %d: the given 3 Cents! is played for free", c)
	}
	if tb.G.Players[0].ExtraLootPlays != 0 {
		t.Error("no loot play is left over")
	}
}
