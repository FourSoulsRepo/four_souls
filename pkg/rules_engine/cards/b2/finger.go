package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Finger! (Passive Treasure Card)
//
//	Each time a player rolls a ❷, you may swap a non-eternal item you control with a non-eternal item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var finger = engine.CardDef{
	Ref:    "finger",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
