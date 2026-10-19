package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Overflow! (Bad Event Card)
//
//	Each player rolls-
//	1-2: They skip their next turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var overflow = engine.CardDef{
	Ref:    "overflow",
	Kind:   engine.EventCard,
	Copies: 1,
}
