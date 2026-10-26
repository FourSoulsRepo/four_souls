package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestSoulOfPride(t *testing.T) {
	tb := bonusTable(t, "soul_of_pride")
	tb.G.AddCountersTo(tb.G.BonusSouls[0], 4)
	tb.G.Players[0].Cents = 20
	tb.G.Turn.Purchases = 2
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "ask.deck.treasure") // 5 counters
	settleAny(tb)
	if n := bonusCounters(tb); n != 5 {
		t.Fatalf("counters = %d, want 5", n)
	}
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "ask.deck.treasure")
	settleAny(tb)
	if s := tb.G.SoulValue(0); s != 1 {
		t.Errorf("Isaac's soul value = %d, want 1: the active player gains it", s)
	}
}
