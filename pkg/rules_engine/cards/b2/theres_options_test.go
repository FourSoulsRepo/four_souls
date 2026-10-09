package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheresOptions(t *testing.T) {
	tb := itemTable(t, items("theres_options"))
	v := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0})
	if v.Players[0].TreasureTop == nil {
		t.Error("the top treasure card is not shown on your turn")
	}
	if v := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 1}); v.Players[0].TreasureTop != nil {
		t.Error("another player sees it")
	}
	if n := tb.G.Turn.Purchases; n != 1 {
		t.Fatal("setup")
	}
}
