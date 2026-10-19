package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Member Card (Active Treasure Card)
//
//	{Tap Effect}Look at the top 3 cards of the treasure deck. Put them back in any order, then put a shop item into discard.
//	Shop items you purchase cost 4¢ less.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var memberCard = engine.CardDef{
	Ref:    "member_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
