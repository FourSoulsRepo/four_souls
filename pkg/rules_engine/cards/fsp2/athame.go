package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Athame (Paid Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Force a player to discard 2 loot cards.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var athame = engine.CardDef{
	Ref:    "athame",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
