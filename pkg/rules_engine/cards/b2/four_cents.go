package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// 4 Cents! (4 Cent Card)
//
//	Gain 4¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fourCents = engine.CardDef{
	Ref:    "four_cents",
	Kind:   engine.LootCard,
	Copies: 9,
}
