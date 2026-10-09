package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Midas Touch (Passive Treasure Card)
//
//	Each time a monster dies, gain 3¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theMidasTouch = engine.CardDef{
	Ref:    "the_midas_touch",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
