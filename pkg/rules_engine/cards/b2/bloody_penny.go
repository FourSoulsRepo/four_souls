package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bloody Penny (Trinket Card)
//
//	Each time a player dies, before paying penalties, loot 1.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bloodyPenny = engine.CardDef{
	Ref:    "bloody_penny",
	Kind:   engine.LootCard,
	Copies: 1,
}
