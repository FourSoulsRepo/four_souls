package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D100 (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1: Loot 1.
//	2: Loot 2.
//	3: Gain 3¢.
//	4: Gain 4¢.
//	5: Gain +1{HP} till end of turn.
//	6: Gain +1{ATK} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theD100 = engine.CardDef{
	Ref:    "the_d100",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
