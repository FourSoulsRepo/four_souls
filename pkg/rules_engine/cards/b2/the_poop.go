package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Poop (Paid Treasure Card)
//
//	Each time you take damage, put a counter on this.
//	{Paid Effect}Remove a counter from this:
//	Prevent the next 1 damage you would take this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var thePoop = engine.CardDef{
	Ref:    "the_poop",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
