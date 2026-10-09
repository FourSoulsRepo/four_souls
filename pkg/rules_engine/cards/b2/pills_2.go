package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Gain 4¢.
//	3-4: Gain 7¢.
//	5-6: Lose 4¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pills2 = engine.CardDef{
	Ref:    "pills_2",
	Kind:   engine.LootCard,
	Copies: 1,
}
