package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D20 (Active Treasure Card)
//
//	{Tap Effect}Reroll an item.
//	(Destroy that item and replace it with the top card of the treasure deck.)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theD20 = engine.CardDef{
	Ref:    "the_d20",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
