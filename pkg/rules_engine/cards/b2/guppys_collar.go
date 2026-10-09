package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Collar (Passive Treasure Card)
//
//	Each time you would die, roll-
//	1-3: Prevent death. If it's your turn, cancel everything that hasn't resolved and end your turn.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guppysCollar = engine.CardDef{
	Ref:    "guppys_collar",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
