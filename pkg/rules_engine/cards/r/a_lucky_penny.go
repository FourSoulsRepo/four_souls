package r

import engine "github.com/FourSoulsRepo/rules_engine"

// A Lucky Penny! (1 Cent Card)
//
//	Gain 1¢.
//	Rewards are doubled till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var aLuckyPenny = engine.CardDef{
	Ref:    "a_lucky_penny",
	Kind:   engine.LootCard,
	Copies: 1,
}
