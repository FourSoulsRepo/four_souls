package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Baleful (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	When you control 1 or 2 souls, you have +1{ATK}. If you control 3+ souls, you have +2{ATK} instead.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theBaleful = engine.CardDef{
	Ref:          "the_baleful",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           1,
	ATK:          1,
	StartingItem: "soulbond",
	ThreePlus:    true,
	Tap:          true,
}
