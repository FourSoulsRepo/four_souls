package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Loot 3. Then put 3 cards from your hand on top of the loot deck in any order.
//	3-4: Loot 3. Then put 1 card from your hand on top of the loot deck.
//	5-6: Loot 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pillsR = engine.CardDef{
	Ref:    "pills_r",
	Kind:   engine.LootCard,
	Copies: 1,
}
