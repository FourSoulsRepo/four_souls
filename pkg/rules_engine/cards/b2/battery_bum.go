package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Battery Bum (Paid Treasure Card)
//
//	{Paid Effect}Pay 4¢:
//	Recharge an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var batteryBum = engine.CardDef{
	Ref:    "battery_bum",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
