package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Compass (Passive Treasure Card)
//
//	At the end of your turn, look at the top 4 cards of the loot deck. Put them back in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theCompass = engine.CardDef{
	Ref:    "the_compass",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
