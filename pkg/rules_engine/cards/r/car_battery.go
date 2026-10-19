package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Car Battery (Paid Treasure Card)
//
//	This can't be recharged except by its own abilities.
//	{Tap Effect}Recharge up to 3 other items in play.
//	{Paid Effect}Pay 7¢: Recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var carBattery = engine.CardDef{
	Ref:    "car_battery",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
