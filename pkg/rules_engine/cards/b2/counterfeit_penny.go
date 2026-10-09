package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Counterfeit Penny (Trinket Card)
//
//	If you would gain any number of ¢, gain that much +1¢ instead.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var counterfeitPenny = engine.CardDef{
	Ref:    "counterfeit_penny",
	Kind:   engine.LootCard,
	Copies: 1,
}
