package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Rock Bottom (Passive Treasure Card)
//
//	When you gain or purchase this, put counters on it equal to the number of loot cards in your hand.
//	At the end of your turn, if you have fewer loot cards in your hand than there are counters on this, loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var rockBottom = engine.CardDef{
	Ref:    "rock_bottom",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
