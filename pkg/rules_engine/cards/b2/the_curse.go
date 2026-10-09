package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Curse (Eternal Treasure Card)
//
//	At the start of your turn, put the top card of a deck into discard.
//	{Tap Effect}Put the top card of any discard on top of its deck.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theCurse = engine.CardDef{
	Ref:     "the_curse",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
