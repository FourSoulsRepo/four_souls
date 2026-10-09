package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Battery (Battery Card)
//
//	Choose a player. Recharge each item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var megaBattery = engine.CardDef{
	Ref:    "mega_battery",
	Kind:   engine.LootCard,
	Copies: 1,
}
