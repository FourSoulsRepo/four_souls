package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// 2 Cents! (2 Cent Card)
//
//	Gain 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var twoCents = engine.CardDef{
	Ref:    "two_cents",
	Kind:   engine.LootCard,
	Copies: 12,
}
