package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mystery Sack (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1-2: Loot 1.
//	3-4: Gain 4¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mysterySack = engine.CardDef{
	Ref:    "mystery_sack",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
