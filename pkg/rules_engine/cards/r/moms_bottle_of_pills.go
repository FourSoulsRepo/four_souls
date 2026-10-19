package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Bottle of Pills (Passive Treasure Card)
//
//	When you gain or purchase this, reveal the top 10 cards of the loot deck. Put each card named Pills! in your hand and the rest on the bottom of the deck.
//	At the start of your turn, put a card named Pills! from the loot discard into your hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsBottleOfPills = engine.CardDef{
	Ref:    "moms_bottle_of_pills",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
