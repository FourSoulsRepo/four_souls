package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Callus (Trinket Card)
//
//	Prevent all non-combat damage you would take.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var callus = engine.CardDef{
	Ref:    "callus",
	Kind:   engine.LootCard,
	Copies: 1,
}
