package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// IX. The Hermit (Wildcard Card)
//
//	Look at the top 5 cards of the treasure deck. Put 1 on top and the rest on the bottom.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ixTheHermit = engine.CardDef{
	Ref:    "ix_the_hermit",
	Kind:   engine.LootCard,
	Copies: 1,
}
