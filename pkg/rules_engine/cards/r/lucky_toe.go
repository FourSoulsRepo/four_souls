package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lucky Toe (Trinket Card)
//
//	If you would loot, except during the loot step, instead loot that much +1.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var luckyToe = engine.CardDef{
	Ref:    "lucky_toe",
	Kind:   engine.LootCard,
	Copies: 1,
}
