package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Spoon Bender (Active Treasure Card)
//
//	{Tap Effect}Add 1 to a roll.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spoonBender = engine.CardDef{
	Ref:    "spoon_bender",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
