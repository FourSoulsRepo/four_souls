package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Shiny Rock (Passive Treasure Card)
//
//	Each time you activate an item, gain 1¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var shinyRock = engine.CardDef{
	Ref:    "shiny_rock",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
