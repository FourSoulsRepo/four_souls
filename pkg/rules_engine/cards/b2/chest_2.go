package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chest (Good Event Card)
//
//	Roll-
//	1-2: Loot 1.
//	3-4: Loot 2.
//	5-6: Loot 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chest2 = engine.CardDef{
	Ref:    "chest_2",
	Kind:   engine.EventCard,
	Copies: 1,
}
