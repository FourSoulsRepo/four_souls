package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D10 (Passive Treasure Card)
//
//	Each time a player rolls a ❸, you may put the top card of the Monster Deck in a monster slot not being attacked.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theD10 = engine.CardDef{
	Ref:    "the_d10",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
