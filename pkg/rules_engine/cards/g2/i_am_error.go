package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// I Am Error! (Good Event Card)
//
//	Roll-
//	1-2: You must make an additional attack on the monster deck.
//	3-4: Loot 3.
//	5-6: Gain +1 Treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iAmError = engine.CardDef{
	Ref:    "i_am_error",
	Kind:   engine.EventCard,
	Copies: 1,
}
