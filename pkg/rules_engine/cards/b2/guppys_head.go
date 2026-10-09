package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Head (Active Treasure Card)
//
//	{Tap Effect}Choose a player. That player gives you a loot card.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guppysHead = engine.CardDef{
	Ref:    "guppys_head",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
