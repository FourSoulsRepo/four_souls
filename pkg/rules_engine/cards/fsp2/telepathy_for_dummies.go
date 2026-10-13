package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Telepathy For Dummies (Passive Treasure Card)
//
//	At the start of your turn, roll-
//	You may change the result of your next roll this turn to this result.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var telepathyForDummies = engine.CardDef{
	Ref:    "telepathy_for_dummies",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
