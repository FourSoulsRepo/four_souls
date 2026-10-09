package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Synthoil (Passive Treasure Card)
//
//	You have +1 to attack rolls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var synthoil = engine.CardDef{
	Ref:    "synthoil",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
