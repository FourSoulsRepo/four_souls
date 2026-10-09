package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Dead Cat (Passive Treasure Card)
//
//	This item starts with 9 counters on it.
//	If you would take damage while this has counters on it, remove that many counters and prevent that much damage.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theDeadCat = engine.CardDef{
	Ref:    "the_dead_cat",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
