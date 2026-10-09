package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Bum (Passive Treasure Card)
//
//	At the start of your turn, roll-
//	1-2: Gain 3¢.
//	3-4: Loot 1.
//	5-6: Take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var darkBum = engine.CardDef{
	Ref:    "dark_bum",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
