package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Glass Cannon (Active Treasure Card)
//
//	{Tap Effect}Destroy another item, then roll-
//	1-5: Destroy this and loot 2.
//	6: Recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var glassCannon = engine.CardDef{
	Ref:    "glass_cannon",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
