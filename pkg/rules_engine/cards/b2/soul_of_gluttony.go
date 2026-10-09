package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Gluttony (Bonus Soul Card)
//
//	The first player to have 10 or more loot cards in their hand gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulOfGluttony = engine.CardDef{
	Ref:    "soul_of_gluttony",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
}
