package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rubber Cement (Passive Treasure Card)
//
//	Each time you miss an attack roll, roll-
//	1-3: Deal 1 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var rubberCement = engine.CardDef{
	Ref:    "rubber_cement",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
