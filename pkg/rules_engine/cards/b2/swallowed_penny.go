package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Swallowed Penny (Trinket Card)
//
//	Each time you take damage, gain 1¢.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var swallowedPenny = engine.CardDef{
	Ref:    "swallowed_penny",
	Kind:   engine.LootCard,
	Copies: 1,
}
