package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// 20/20 (Active Treasure Card)
//
//	{Tap Effect}Add up to 2 to an attack roll.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var card2020 = engine.CardDef{
	Ref:    "20_20",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
