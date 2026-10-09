package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Battery (Active Treasure Card)
//
//	{Tap Effect}Recharge another item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theBattery = engine.CardDef{
	Ref:    "the_battery",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
