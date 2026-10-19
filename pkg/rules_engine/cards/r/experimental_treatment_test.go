package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestExperimentalTreatmentSoul(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "experimental_treatment"), seat("cain")}})
	tb.G.ForceRolls(6)
	tb.Activate(0, "experimental_treatment", 0)
	if tb.G.SoulValue(0) != 1 {
		t.Fatal("6: it becomes his soul")
	}
	// As a soul it keeps its ability (R-CARD-38, S-YUG).
	tb.G.ForceRolls(3)
	tb.Activate(0, "experimental_treatment", 0)
	if tb.G.SoulValue(0) != 0 || tb.G.Players[0].Cents != 9 {
		t.Errorf("the soul is destroyed for another roll: souls %d, cents %d", tb.G.SoulValue(0), tb.G.Players[0].Cents)
	}
}

func TestExperimentalTreatmentDie(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "experimental_treatment", "spoon_bender"), seat("cain")}})
	tb.G.ForceRolls(1)
	tb.Activate(0, "experimental_treatment", 0, "spoon_bender")
	if !tb.G.Players[0].Dead {
		t.Error("1: you die")
	}
}

func TestExperimentalTreatmentCancelled(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "experimental_treatment"), {Character: "cain", Hand: []engine.CardRef{"butter_bean"}},
	}})
	tb.Start(0, "experimental_treatment", 0)
	tb.G.Players[1].ExtraLootPlays = 1
	tb.Pass(0)
	tb.Play(1, "butter_bean", "experimental_treatment: $ Destroy this: roll- 1: You die. 2: Loot 3. 3: Gain 9¢. 4: Gain +1 treasure. 5: Gain +2 treasure. 6: This becomes a soul. Gain it.")
	if hasItem(tb.G, 0, "experimental_treatment") {
		t.Error("destroying it is the cost: it stays destroyed (S-YUG)")
	}
}
