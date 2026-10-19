package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Wishbone (Trinket Card)
//
//	When this is destroyed, gain +1 treasure.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var wishbone = engine.CardDef{
	Ref:    "wishbone",
	Kind:   engine.LootCard,
	Copies: 1,
}
