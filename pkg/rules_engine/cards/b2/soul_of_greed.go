package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Greed (Bonus Soul Card)
//
//	The first player to have 25¢ or more gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulOfGreed = engine.CardDef{
	Ref:    "soul_of_greed",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
}
