package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Trick Penny (Passive Treasure Card)
//
//	Shop items you purchase cost 1¢ less.
//	$ items you control cost 1¢ less to activate.
//	Each time you purchase from the shop or treasure deck, gain 1¢.
//	At the end of your turn, gain 1¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var trickPenny = engine.CardDef{
	Ref:    "trick_penny",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
