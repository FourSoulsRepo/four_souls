package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Abaddon (Passive Treasure Card)
//
//	+3{ATK}
//	Each time you take damage, die!
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var abaddon = engine.CardDef{
	Ref:    "abaddon",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
