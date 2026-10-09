package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Penny! (1 Cent Card)
//
//	Gain 1¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var aPenny6 = engine.CardDef{
	Ref:    "a_penny_6",
	Kind:   engine.LootCard,
	Copies: 6,
}
