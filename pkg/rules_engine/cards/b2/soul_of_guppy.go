package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Guppy (Bonus Soul Card)
//
//	The first player to control 2 or more guppy items gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulOfGuppy = engine.CardDef{
	Ref:    "soul_of_guppy",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
}
