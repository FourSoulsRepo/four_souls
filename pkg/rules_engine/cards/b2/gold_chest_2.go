package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gold Chest (Good Event Card)
//
//	Roll-
//	1-2: Gain +1 Treasure.
//	3-4: Loot 1.
//	5-6: Loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goldChest2 = engine.CardDef{
	Ref:    "gold_chest_2",
	Kind:   engine.EventCard,
	Copies: 1,
}
