package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bag Lunch (Trinket Card)
//
//	If this would be destroyed, if it has no counters on it, put a counter on it instead.
//	You have +1{HP} while this has a counter on it.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bagLunch = engine.CardDef{
	Ref:    "bag_lunch",
	Kind:   engine.LootCard,
	Copies: 1,
}
