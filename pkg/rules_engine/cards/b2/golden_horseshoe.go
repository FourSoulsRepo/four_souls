package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Horseshoe (Trinket Card)
//
//	At the start of your turn, look at the top card of the treasure deck, you may put it on the bottom.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goldenHorseshoe = engine.CardDef{
	Ref:    "golden_horseshoe",
	Kind:   engine.LootCard,
	Copies: 1,
}
