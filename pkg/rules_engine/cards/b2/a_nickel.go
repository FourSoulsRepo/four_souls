package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Nickel! (Nickel Card)
//
//	Gain 5¢
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var aNickel = engine.CardDef{
	Ref:    "a_nickel",
	Kind:   engine.LootCard,
	Copies: 5,
}
