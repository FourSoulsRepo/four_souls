package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Flush! (Active Treasure Card)
//
//	{Tap Effect}Choose one- Put each monster not being attacked on the bottom of the monster deck. Put each shop item on the bottom of the treasure deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var flush = engine.CardDef{
	Ref:    "flush",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
