package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// O. The Fool (Wildcard Card)
//
//	End the turn. Cancel everything that hasn't resolved.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var oTheFool = engine.CardDef{
	Ref:    "o_the_fool",
	Kind:   engine.LootCard,
	Copies: 1,
}
