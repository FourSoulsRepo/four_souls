package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestHolyMomsEye(t *testing.T) {
	tb := slayTable(t, "holy_moms_eye", seat("isaac", "mystery_sack"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "mystery_sack", 0, "mystery_sack")
	if !tb.G.Object(tb.Find(0, "mystery_sack")).Charged {
		t.Error("the sack was not recharged")
	}
}

// TestHolyMomsEyeOffersOnlySpentActiveItems: the recharge question
// leaves out passive items and items that are already charged (19.145).
func TestHolyMomsEyeOffersOnlySpentActiveItems(t *testing.T) {
	cain := engine.SituationPlayer{Character: "cain", Items: items("jawbone", "breakfast"), Deactivated: items("jawbone")}
	tb := slayTable(t, "holy_moms_eye", seat("isaac", "mystery_sack", "the_d6"), cain)
	tb.G.ForceRolls(2)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(0, "mystery_sack")}}); err != nil {
		t.Fatal(err)
	}
	tb.UntilChoose()
	want := []string{"mystery_sack", "jawbone", "ask.none"}
	if got := tb.Options(); !slices.Equal(got, want) {
		t.Errorf("options = %v, want %v", got, want)
	}
}

// TestChargedBabyOffersOnlySpentActiveItems: Charged Baby uses the same
// recharge question (19.145).
func TestChargedBabyOffersOnlySpentActiveItems(t *testing.T) {
	tb := itemTable(t, items("charged_baby", "mystery_sack", "the_d6", "breakfast"))
	tb.G.ForceRolls(2)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(0, "mystery_sack")}}); err != nil {
		t.Fatal(err)
	}
	tb.UntilChoose()
	want := []string{"mystery_sack", "ask.none"}
	if got := tb.Options(); !slices.Equal(got, want) {
		t.Errorf("options = %v, want %v", got, want)
	}
}
