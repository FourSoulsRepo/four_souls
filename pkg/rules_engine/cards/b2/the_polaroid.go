package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Polaroid (Passive Treasure Card)
//
//	At the end of your turn, if you have 0 loot cards in your hand, loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var thePolaroid = engine.CardDef{
	Ref:    "the_polaroid",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
