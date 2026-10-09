package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XVII. The Stars (Wildcard Card)
//
//	Gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xviiTheStars = engine.CardDef{
	Ref:    "xvii_the_stars",
	Kind:   engine.LootCard,
	Copies: 1,
}
