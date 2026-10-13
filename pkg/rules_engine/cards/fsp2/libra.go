package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Libra (Active Treasure Card)
//
//	{Tap Effect}Change the result of a dice roll to a 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var libra = engine.CardDef{
	Ref:    "libra",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
