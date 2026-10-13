package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Smart Fly (Active Treasure Card)
//
//	{Tap Effect}Look at the top card of a deck. You may put it into discard or put it back on top.
//	Each time you take damage, you may recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var smartFly = engine.CardDef{
	Ref:    "smart_fly",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
