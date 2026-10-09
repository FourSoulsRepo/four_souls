package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// I. The Magician (Wildcard Card)
//
//	Change the result of a dice roll to a number of your choosing.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iTheMagician = engine.CardDef{
	Ref:    "i_the_magician",
	Kind:   engine.LootCard,
	Copies: 1,
}
