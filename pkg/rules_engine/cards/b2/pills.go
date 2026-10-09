package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Loot 1.
//	3-4: Loot 3.
//	5-6: Discard 1 loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pills = engine.CardDef{
	Ref:    "pills",
	Kind:   engine.LootCard,
	Copies: 1,
}
