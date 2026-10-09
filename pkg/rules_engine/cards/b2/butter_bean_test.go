package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestButterBean(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"a_dime"}},
		{Character: "cain", Hand: []engine.CardRef{"butter_bean"}},
	}}, Set)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "a_dime")}}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.G.Players[1].ExtraLootPlays = 1 // e.g. from Cain's ↷ (R-CARD-08)
	tb.Play(1, "butter_bean", "a_dime")
	if c := tb.G.Players[0].Cents; c != 0 {
		t.Errorf("the cancelled A Dime!! gave %d¢", c)
	}
	if !inDiscard(tb.G, "a_dime") || !inDiscard(tb.G, "butter_bean") {
		t.Error("both loot cards should be in the discard")
	}
}
