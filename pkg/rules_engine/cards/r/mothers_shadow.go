package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mother’s Shadow! (Bad Event Card)
//
//	The player who controls the most items or tied for the most dies. Put this in the monster deck 6 cards from the top.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mothersShadow = engine.CardDef{
	Ref:    "mothers_shadow",
	Kind:   engine.EventCard,
	Copies: 1,
}
