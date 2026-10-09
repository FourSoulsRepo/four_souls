package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ehwaz (Pill/Rune Card)
//
//	Put each monster not being attacked into discard and replace each with the top card of the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ehwaz = engine.CardDef{
	Ref:    "ehwaz",
	Kind:   engine.LootCard,
	Copies: 1,
}
