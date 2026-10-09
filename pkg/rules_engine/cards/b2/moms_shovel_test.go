package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestMomsShovel(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, {Character: "cain", Souls: items("monstro")}},
		Shop:    items("moms_shovel"),
	}, Set)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "moms_shovel")
	if tb.G.Object(tb.Find(0, "moms_shovel")).Charged {
		t.Fatal("the shovel entered play charged")
	}
	tb.EndTurn()
	tb.EndTurn() // back to Isaac: recharged
	tb.Activate(0, "moms_shovel", 0, "monstro")
	if tb.G.SoulValue(0) != 1 || tb.G.SoulValue(1) != 0 {
		t.Error("the soul was not stolen")
	}
}
