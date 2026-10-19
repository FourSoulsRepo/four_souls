package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Judas’ Shadow (Passive Treasure Card)
//
//	Each time you die, before paying penalties, put 2 counters on this.
//	When this is destroyed, gain x¢ and loot x, where x is equal to the number of counters on this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var judasShadow = engine.CardDef{
	Ref:    "judas_shadow",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
