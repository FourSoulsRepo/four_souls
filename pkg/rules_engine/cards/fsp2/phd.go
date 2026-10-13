package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// P.H.D. (Passive Treasure Card)
//
//	You may add or subtract 1 from any of your non-attack rolls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var phd = engine.CardDef{
	Ref:    "phd",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
