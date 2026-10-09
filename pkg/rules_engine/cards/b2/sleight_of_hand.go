package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sleight Of Hand (Eternal Treasure Card)
//
//	{Tap Effect}Look at the top 5 cards of a deck. Put them back in any order.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sleightOfHand = engine.CardDef{
	Ref:     "sleight_of_hand",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
