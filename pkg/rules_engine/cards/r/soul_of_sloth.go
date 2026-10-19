package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Sloth (Bonus Soul Card)
//
//	The first time a player controls 4 items, the active player chooses a player who controls the fewest items or tied for fewest. That player gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulOfSloth = engine.CardDef{
	Ref:    "soul_of_sloth",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
}
