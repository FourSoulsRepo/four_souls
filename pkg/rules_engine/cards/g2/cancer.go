package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cancer (Trinket Card)
//
//	You have +1 to your first dice roll each turn.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cancer = engine.CardDef{
	Ref:    "cancer",
	Kind:   engine.LootCard,
	Copies: 1,
}
