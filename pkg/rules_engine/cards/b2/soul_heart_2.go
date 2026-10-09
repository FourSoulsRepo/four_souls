package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Heart (Dice Shard/Soul Heart Card)
//
//	Choose a player. Prevent the next 1 damage they would take this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var soulHeart2 = engine.CardDef{
	Ref:    "soul_heart_2",
	Kind:   engine.LootCard,
	Copies: 2,
}
