package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ace Of Diamonds (Wildcard Card)
//
//	Gain x¢, where x is the number of monster slots plus the number of loot cards in your hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var aceOfDiamonds = engine.CardDef{
	Ref:    "ace_of_diamonds",
	Kind:   engine.LootCard,
	Copies: 1,
}
