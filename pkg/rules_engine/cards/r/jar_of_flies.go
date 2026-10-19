package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Jar Of Flies (Active Treasure Card)
//
//	Each time a monster dies, put a counter on this.
//	{Tap Effect}Remove 1 or more counters from this:
//	Loot 1. If 3+ counters were removed, deal 1 damage to a monster instead.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var jarOfFlies = engine.CardDef{
	Ref:    "jar_of_flies",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
