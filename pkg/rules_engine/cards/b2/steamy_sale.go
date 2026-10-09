package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Steamy Sale! (Passive Treasure Card)
//
//	Shop items you purchase cost 5¢ less.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var steamySale = engine.CardDef{
	Ref:    "steamy_sale",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
