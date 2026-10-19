package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lunch (Passive Treasure Card)
//
//	+1{HP}
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lunch = engine.CardDef{
	Ref:    "lunch",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
