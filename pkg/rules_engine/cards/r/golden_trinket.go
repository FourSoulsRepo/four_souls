package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Trinket (Trinket Card)
//
//	When this enters play put a gold counter on another non-eternal item you control.
//	This has the abilities of other items with gold counters on them.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goldenTrinket = engine.CardDef{
	Ref:    "golden_trinket",
	Kind:   engine.LootCard,
	Copies: 1,
}
