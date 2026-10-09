package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Godhead (Active Treasure Card)
//
//	{Tap Effect}Change the result of a dice roll to a 1 or 6.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var godhead = engine.CardDef{
	Ref:    "godhead",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
