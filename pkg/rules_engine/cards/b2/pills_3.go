package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: You gain +1{ATK} till the end of turn.
//	3-4: You gain +1{HP} till the end of turn.
//	5-6: Take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pills3 = engine.CardDef{
	Ref:    "pills_3",
	Kind:   engine.LootCard,
	Copies: 1,
}
