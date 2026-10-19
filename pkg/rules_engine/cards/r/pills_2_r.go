package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll 3 times!-
//	1-2: You gain +1{ATK} till end of turn!
//	3-4: You gain +1{HP} till end of turn!
//	5-6: You take 1 damage!
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pills2R = engine.CardDef{
	Ref:    "pills_2_r",
	Kind:   engine.LootCard,
	Copies: 1,
}
