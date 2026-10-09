package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XV. The Devil (Wildcard Card)
//
//	Destroy an item you control. If you do, steal a non-eternal item from a player or from the shop.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xvTheDevil = engine.CardDef{
	Ref:    "xv_the_devil",
	Kind:   engine.LootCard,
	Copies: 1,
}
