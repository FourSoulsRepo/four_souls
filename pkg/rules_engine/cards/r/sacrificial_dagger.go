package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sacrificial Dagger (Paid Treasure Card)
//
//	Each time a player rolls a ❶, put a counter on this.
//	{Paid Effect}Remove 3 counters from this:
//	Deal 1 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sacrificialDagger = engine.CardDef{
	Ref:    "sacrificial_dagger",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
