package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Eye (Passive Treasure Card)
//
//	Each other player plays with their hand revealed.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guppysEye = engine.CardDef{
	Ref:    "guppys_eye",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Guppy:  true,
}
