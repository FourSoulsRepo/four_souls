package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Habit (Passive Treasure Card)
//
//	The first time you take damage each turn, you may recharge an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theHabit = engine.CardDef{
	Ref:    "the_habit",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
