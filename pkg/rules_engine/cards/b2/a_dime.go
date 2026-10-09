package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Dime!! (Nickel Card)
//
//	Gain 10¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var aDime = engine.CardDef{
	Ref:    "a_dime",
	Kind:   engine.LootCard,
	Copies: 1,
}
