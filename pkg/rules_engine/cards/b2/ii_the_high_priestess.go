package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// II. The High Priestess (Wildcard Card)
//
//	Choose a player or monster, then roll-
//	Deal damage to them equal to the result.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iiTheHighPriestess = engine.CardDef{
	Ref:    "ii_the_high_priestess",
	Kind:   engine.LootCard,
	Copies: 1,
}
