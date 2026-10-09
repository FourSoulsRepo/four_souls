package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blank Rune (Pill/Rune Card)
//
//	Roll-
//	1: Each player gains 1¢.
//	2: Each player loots 2.
//	3: Each player takes 3 damage.
//	4: Each player gains 4¢.
//	5: Each player loots 5.
//	6: Each player gains 6¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blankRune = engine.CardDef{
	Ref:    "blank_rune",
	Kind:   engine.LootCard,
	Copies: 1,
}
