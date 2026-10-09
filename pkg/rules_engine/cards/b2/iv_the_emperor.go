package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// IV. The Emperor (Wildcard Card)
//
//	Look at the top 5 cards of the monster deck. Put 1 on top and the rest on the bottom.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ivTheEmperor = engine.CardDef{
	Ref:    "iv_the_emperor",
	Kind:   engine.LootCard,
	Copies: 1,
}
