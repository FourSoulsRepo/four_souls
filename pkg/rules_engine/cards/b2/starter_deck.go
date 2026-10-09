package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Starter Deck (Passive Treasure Card)
//
//	At the end of your turn, if you have 8 or more loot cards in your hand, loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var starterDeck = engine.CardDef{
	Ref:    "starter_deck",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
