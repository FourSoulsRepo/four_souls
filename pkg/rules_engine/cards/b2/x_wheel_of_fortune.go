package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// X. Wheel Of Fortune (Wildcard Card)
//
//	Roll-
//	1: Gain 1¢.
//	2: Take 2 damage.
//	3. Loot 3.
//	4. Lose 4¢.
//	5: Gain 5¢.
//	6: Gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xWheelOfFortune = engine.CardDef{
	Ref:    "x_wheel_of_fortune",
	Kind:   engine.LootCard,
	Copies: 1,
}
