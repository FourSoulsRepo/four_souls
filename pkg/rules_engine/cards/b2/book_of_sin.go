package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Book Of Sin (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1-2: Gain 1¢.
//	3-4: Loot 1.
//	5-6: Gain +1{HP} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bookOfSin = engine.CardDef{
	Ref:    "book_of_sin",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
