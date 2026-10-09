package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chest (Good Event Card)
//
//	Roll-
//	1-2: Gain 1¢.
//	3-4: Gain 3¢.
//	5-6: Gain 6¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var chest = engine.CardDef{
	Ref:    "chest",
	Kind:   engine.EventCard,
	Copies: 1,
}
