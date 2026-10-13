package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lard (Passive Treasure Card)
//
//	+2{HP}
//	Each time you take damage, discard a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lard = engine.CardDef{
	Ref:    "lard",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
