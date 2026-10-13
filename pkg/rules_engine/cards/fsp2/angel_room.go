package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Angel Room (Good Event Card)
//
//	Roll-
//	1: Gain +2 treasure.
//	2-3: Gain +1 treasure.
//	4-6: Loot 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var angelRoom = engine.CardDef{
	Ref:    "angel_room",
	Kind:   engine.EventCard,
	Copies: 1,
}
