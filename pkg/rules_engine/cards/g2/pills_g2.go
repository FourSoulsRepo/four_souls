package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Recharge each item you control.
//	3-4: Monsters you attack have -1{DC} till end of turn.
//	5-6: Monsters you attack have +2{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pillsG2 = engine.CardDef{
	Ref:    "pills_g2",
	Kind:   engine.LootCard,
	Copies: 1,
}
