package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dead Weight (Eternal Treasure Card)
//
//	You take no combat damage on attack rolls of 1.
//	{Tap Effect}Cancel the triggered ability of a monster or non-eternal item.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var deadWeight = engine.CardDef{
	Ref:     "dead_weight",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
