package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spelunker Hat (Passive Treasure Card)
//
//	At the start of your turn, you may look at the top 3 cards of a deck and put them back in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spelunkerHat = engine.CardDef{
	Ref:    "spelunker_hat",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
