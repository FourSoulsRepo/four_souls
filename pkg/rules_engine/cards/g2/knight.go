package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Knight (Basic Monster Card)
//
//	The attacking player also deals combat damage to this on attack rolls of 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var knight = engine.CardDef{
	Ref:     "knight",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      6,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
