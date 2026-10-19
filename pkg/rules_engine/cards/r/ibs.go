package r

import engine "github.com/FourSoulsRepo/rules_engine"

// IBS (Eternal Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Choose one- Cancel your attack on a monster. Gain 3¢.
//	{Paid Effect}Discard 2 loot cards:
//	Cancel the ↷ or $ ability of an item or loot being played.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ibs = engine.CardDef{
	Ref:     "ibs",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
