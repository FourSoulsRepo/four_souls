package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Stays Go (24.23): Butter Bean cancels the ability on the stack, but
// its cost destroyed the card, and setup.stack can't hold an ability
// whose card is no longer in play.
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
