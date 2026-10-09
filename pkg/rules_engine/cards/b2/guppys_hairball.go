package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Hairball (Trinket Card)
//
//	Each time you would take damage, roll-
//	6: Prevent 1 of that damage.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var guppysHairball = engine.CardDef{
	Ref:    "guppys_hairball",
	Kind:   engine.LootCard,
	Copies: 1,
}
