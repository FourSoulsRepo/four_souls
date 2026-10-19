package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Two Of Spades (Wildcard Card)
//
//	Choose another player. Loot x, where x is the number of loot cards in that player's hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var twoOfSpades = engine.CardDef{
	Ref:    "two_of_spades",
	Kind:   engine.LootCard,
	Copies: 1,
}
