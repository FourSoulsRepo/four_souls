package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Jera (Pill/Rune Card)
//
//	Loot x, where x is equal to the number of loot cards in your hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var jera = engine.CardDef{
	Ref:    "jera",
	Kind:   engine.LootCard,
	Copies: 1,
}
