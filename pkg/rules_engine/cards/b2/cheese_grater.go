package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cheese Grater (Passive Treasure Card)
//
//	Each time a player rolls a ❻, reveal the top card of any deck. Put it back or put it into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cheeseGrater = engine.CardDef{
	Ref:    "cheese_grater",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
