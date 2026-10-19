package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Undefined (Active Treasure Card)
//
//	{Tap Effect}Choose a shop item. This becomes a copy of that item till the start of your next turn. Recharge this.
//	At the start of your turn you may put a shop item into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var undefined = engine.CardDef{
	Ref:    "undefined",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
