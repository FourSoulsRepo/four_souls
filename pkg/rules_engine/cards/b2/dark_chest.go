package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Chest (Good Event Card)
//
//	Roll-
//	1-2: Loot 1.
//	3-4: Gain 3¢.
//	5-6: Take 2 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var darkChest = engine.CardDef{
	Ref:    "dark_chest",
	Kind:   engine.EventCard,
	Copies: 1,
}
