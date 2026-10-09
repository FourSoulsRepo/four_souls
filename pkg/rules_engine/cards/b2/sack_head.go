package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sack Head (Active Treasure Card)
//
//	{Tap Effect}Look at the top card of a deck. You may put that card on the bottom of that deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sackHead = engine.CardDef{
	Ref:    "sack_head",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
