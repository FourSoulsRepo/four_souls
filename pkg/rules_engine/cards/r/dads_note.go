package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dad’s Note (Wildcard Card)
//
//	Choose up to 3 non-event monster cards in discard. Put them in one or more monster slots not being attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dadsNote = engine.CardDef{
	Ref:    "dads_note",
	Kind:   engine.LootCard,
	Copies: 1,
}
