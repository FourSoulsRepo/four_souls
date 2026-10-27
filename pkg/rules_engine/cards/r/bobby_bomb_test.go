package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBobbyBomb(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}})
	tb.G.PutOnDeck(engine.LootDeck, "mega_bomb", "a_nickel", "bomb")
	tb.G.PutOnDeck(engine.TreasureDeck, "bobby_bomb")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "ask.deck.treasure")
	if !handHas(tb.G, 0, "bomb") || !handHas(tb.G, 0, "mega_bomb") || handHas(tb.G, 0, "a_nickel") {
		t.Error("the bombs of the top 10 go to his hand; the nickel to the bottom")
	}
}

func TestBobbyBombStartOfTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "bobby_bomb")}})
	tb.G.PutOnDeck(engine.LootDeck, "troll_bombs")
	tb.G.MillTop(engine.LootDeck)
	tb.EndTurn("troll_bombs")
	if !handHas(tb.G, 1, "troll_bombs") {
		t.Error("a bomb card comes back from the loot discard")
	}
}
