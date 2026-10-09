package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Coin Purse (Passive Treasure Card)
//
//	Loot +1 during your loot step.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsCoinPurse = engine.CardDef{
	Ref:    "moms_coin_purse",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
