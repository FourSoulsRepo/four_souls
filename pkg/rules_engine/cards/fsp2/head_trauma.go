package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Head Trauma (Bad Event Card)
//
//	Discard 2 loot cards.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var headTrauma = engine.CardDef{
	Ref:    "head_trauma",
	Kind:   engine.EventCard,
	Copies: 1,
}
