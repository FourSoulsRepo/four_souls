package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCurseOfLoss(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: items("xiii_death")}, {Character: "cain", Souls: items("monstro")}},
		Monsters: items("fly"),
	}, Set)
	tb.RevealFromDeck("curse_of_loss", foe)
	tb.Play(0, "xiii_death", foe, "monstro")
	if s := tb.G.SoulValue(1); s != 0 {
		t.Errorf("Cain's souls = %d, want 0", s)
	}
	if hasCurse(tb.G, 1, "curse_of_loss") {
		t.Error("the curse stays after death")
	}
}
