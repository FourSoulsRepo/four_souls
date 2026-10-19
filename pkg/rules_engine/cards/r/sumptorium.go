package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sumptorium (Eternal Treasure Card)
//
//	Each time you take damage, put a counter on this.
//	{Paid Effect}Remove 2 counters from this:
//	A monster gains -1{DC} till end of turn.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sumptorium = engine.CardDef{
	Ref:     "sumptorium",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
