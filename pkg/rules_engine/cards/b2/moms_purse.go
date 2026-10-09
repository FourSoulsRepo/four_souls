package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Purse (Passive Treasure Card)
//
//	Loot +1 during your loot step.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsPurse = engine.CardDef{
	Ref:    "moms_purse",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
