package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Rib Of Greed (Trinket Card)
//
//	Abilities and the death penalty can't make you discard loot cards or lose ¢.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ribOfGreed = engine.CardDef{
	Ref:    "rib_of_greed",
	Kind:   engine.LootCard,
	Copies: 1,
}
