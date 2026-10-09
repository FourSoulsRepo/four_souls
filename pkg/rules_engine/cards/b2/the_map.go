package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Map (Passive Treasure Card)
//
//	At the end of your turn, look at the top 4 cards of the monster deck. Put them back in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theMap = engine.CardDef{
	Ref:    "the_map",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
