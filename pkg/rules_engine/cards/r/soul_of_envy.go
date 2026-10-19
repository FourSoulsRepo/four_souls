package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Envy (Bonus Soul Card)
//
//	The first time a player controls their 3rd soul, the active player chooses a player who controls the fewest souls or tied for fewest. That player gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulOfEnvy = engine.CardDef{
	Ref:    "soul_of_envy",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
}
