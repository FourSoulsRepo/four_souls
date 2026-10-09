package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// 3 Cents! (3 Cent Card)
//
//	Gain 3¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var threeCents = engine.CardDef{
	Ref:    "three_cents",
	Kind:   engine.LootCard,
	Copies: 15,
}
