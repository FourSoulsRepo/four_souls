package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hagalaz (Pill/Rune Card)
//
//	Look at the top 5 cards of the monster or room deck and put them back in any order, then put a room or monster not being attacked into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hagalaz = engine.CardDef{
	Ref:    "hagalaz",
	Kind:   engine.LootCard,
	Copies: 1,
}
