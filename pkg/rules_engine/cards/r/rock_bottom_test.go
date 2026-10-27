package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestRockBottom(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Cents: 10, Hand: []engine.CardRef{"bomb", "bomb"}}, seat("cain"),
	}})
	tb.G.PutOnDeck(engine.TreasureDeck, "rock_bottom")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "ask.deck.treasure")
	if got := tb.G.Object(tb.Find(0, "rock_bottom")).CountersOf(""); got != 2 {
		t.Fatalf("counters = %d, want his 2 cards", got)
	}
	tb.G.Players[0].Hand = tb.G.Players[0].Hand[:1]
	tb.EndTurn()
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 1 + 2: fewer cards than counters", got)
	}
}
