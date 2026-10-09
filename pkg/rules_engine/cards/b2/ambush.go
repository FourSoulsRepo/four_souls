package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ambush! (Bad Event Card)
//
//	The active player must attack the monster deck 2 times this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ambush = engine.CardDef{
	Ref:    "ambush",
	Kind:   engine.EventCard,
	Copies: 1,
}
