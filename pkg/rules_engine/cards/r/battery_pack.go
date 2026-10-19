package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Battery Pack (Paid Treasure Card)
//
//	At the end of your turn, put a counter on this.
//	{Paid Effect}Remove 2 counters from this:
//	Recharge each item a player controls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var batteryPack = engine.CardDef{
	Ref:    "battery_pack",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
