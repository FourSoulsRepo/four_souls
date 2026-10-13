package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Heart! (Epic Boss Card)
//
//	When a player gains this soul, choose a player who controls the most souls or tied for the most. That player wins.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsHeart = engine.CardDef{
	Ref:    "moms_heart",
	Kind:   engine.MonsterCard,
	Copies: 1,
	HP:     8,
	DC:     4,
	ATK:    2,
	Soul:   2,
}
