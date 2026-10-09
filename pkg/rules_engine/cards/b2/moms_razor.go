package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Razor (Passive Treasure Card)
//
//	Each time a player rolls a ❻, you may deal 1 damage to them.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsRazor = engine.CardDef{
	Ref:    "moms_razor",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
