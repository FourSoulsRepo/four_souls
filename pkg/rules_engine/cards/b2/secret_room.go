package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Secret Room! (Good Event Card)
//
//	Roll-
//	1: Take 3 damage.
//	2-3: Discard 2 loot cards.
//	4-5: Gain 7¢.
//	6: Gain +1 Treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var secretRoom = engine.CardDef{
	Ref:    "secret_room",
	Kind:   engine.EventCard,
	Copies: 1,
}
