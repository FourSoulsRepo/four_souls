package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed! (Bad Event Card)
//
//	Choose the player with the most ¢ or tied for the most. That player loses all their ¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var greedEvent = engine.CardDef{
	Ref:    "greed_event",
	Kind:   engine.EventCard,
	Copies: 1,
}
