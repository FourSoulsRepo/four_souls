package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fruitcake (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1: Loot 1.
//	2-3: Gain 1¢.
//	4-6: Look at the top card of a deck.
//	Each time a player rolls the same result twice in a row on the same turn, recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fruitcake = engine.CardDef{
	Ref:    "fruitcake",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
