package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Sack (Wildcard Card)
//
//	Loot 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var aSack = engine.CardDef{
	Ref:    "a_sack",
	Kind:   engine.LootCard,
	Copies: 1,
}
