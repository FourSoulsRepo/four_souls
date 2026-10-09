package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lil Battery (Battery Card)
//
//	Recharge an item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lilBattery = engine.CardDef{
	Ref:    "lil_battery",
	Kind:   engine.LootCard,
	Copies: 4,
}
