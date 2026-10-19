package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hypercoagulation (Eternal Treasure Card)
//
//	At the start of your turn put a counter on this.
//	You have +1{HP} for each counter on this.
//	At the end of your turn, if this has 3+ counters, remove all of them and loot 3.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hypercoagulation = engine.CardDef{
	Ref:     "hypercoagulation",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
}
