package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dad’s Lost Coin (Passive Treasure Card)
//
//	Each time a player would roll a ❶, you may force that player to reroll it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dadsLostCoin = engine.CardDef{
	Ref:    "dads_lost_coin",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
