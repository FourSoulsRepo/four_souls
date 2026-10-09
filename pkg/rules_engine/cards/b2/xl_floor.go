package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Xl Floor! (Good Event Card)
//
//	Expand monster slots by 1.
//	The active player may attack an additional time this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xlFloor = engine.CardDef{
	Ref:    "xl_floor",
	Kind:   engine.EventCard,
	Copies: 1,
}
