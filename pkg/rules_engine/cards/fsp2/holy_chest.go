package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Chest (Good Event Card)
//
//	Roll-
//	1-2: Prevent the next 2 damage you would take this turn. You may attack an additional time this turn.
//	3-5: Gain 7¢.
//	6: This becomes a soul. Gain it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyChest = engine.CardDef{
	Ref:    "holy_chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Soul:   1,
}
