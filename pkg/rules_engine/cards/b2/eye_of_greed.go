package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Eye Of Greed (Passive Treasure Card)
//
//	Each time a player rolls a ❺, gain 3¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var eyeOfGreed = engine.CardDef{
	Ref:    "eye_of_greed",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
