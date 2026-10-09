package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Restock (Passive Treasure Card)
//
//	At the start of your turn, you may put any number of shop items into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var restock = engine.CardDef{
	Ref:    "restock",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
