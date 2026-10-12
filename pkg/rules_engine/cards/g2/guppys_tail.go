package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Tail (Passive Treasure Card)
//
//	Loot +1 during your loot step.
//	At the end of your turn, discard a loot card.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guppysTail = engine.CardDef{
	Ref:    "guppys_tail",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Guppy:  true,
}
