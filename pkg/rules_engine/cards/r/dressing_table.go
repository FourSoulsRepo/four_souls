package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dressing Table (Good Event Card)
//
//	Reroll any number of items you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dressingTable = engine.CardDef{
	Ref:    "dressing_table",
	Kind:   engine.EventCard,
	Copies: 1,
}
