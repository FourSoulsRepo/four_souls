package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Heart (Dice Shard/Soul Heart Card)
//
//	Choose a player. Prevent the next 1 damage they would take this turn. Till end of turn, when that player dies, deal 2 damage to each player other than that player and you.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var blackHeart = engine.CardDef{
	Ref:    "black_heart",
	Kind:   engine.LootCard,
	Copies: 2,
}
