package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dinner (Passive Treasure Card)
//
//	{HP}
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dinner = engine.CardDef{
	Ref:    "dinner",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
