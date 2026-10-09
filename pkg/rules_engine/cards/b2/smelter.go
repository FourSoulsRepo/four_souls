package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Smelter (Paid Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Gain 3¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var smelter = engine.CardDef{
	Ref:    "smelter",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
