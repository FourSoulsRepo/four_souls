package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Troll Bombs (Bad Event Card)
//
//	Take 2 damage!
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var trollBombs = engine.CardDef{
	Ref:    "troll_bombs",
	Kind:   engine.EventCard,
	Copies: 1,
}
