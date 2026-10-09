package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cambion Conception (Passive Treasure Card)
//
//	Each time you take damage, put counters on this equal to the amount of damage taken. Then, if this has 6+ counters, remove 6 counters from this and gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cambionConception = engine.CardDef{
	Ref:    "cambion_conception",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
