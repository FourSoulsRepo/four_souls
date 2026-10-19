package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bag Of Crafting (Eternal Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Put a counter on this.
//	{Paid Effect}Remove 4 counters from this:
//	Gain +1 treasure.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bagOfCrafting = engine.CardDef{
	Ref:     "bag_of_crafting",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
