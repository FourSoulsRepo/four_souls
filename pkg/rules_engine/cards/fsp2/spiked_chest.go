package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Spiked Chest (Good Event Card)
//
//	Roll-
//	1-2: Take 1 damage.
//	3-4: Take 1 damage. Loot 1.
//	5-6: Take 1 damage. Gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spikedChest = engine.CardDef{
	Ref:    "spiked_chest",
	Kind:   engine.EventCard,
	Copies: 1,
}
