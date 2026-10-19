package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pitfall! (Bad Event Card)
//
//	End the turn. Cancel everything that hasn't resolved.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pitfall = engine.CardDef{
	Ref:    "pitfall",
	Kind:   engine.EventCard,
	Copies: 1,
}
