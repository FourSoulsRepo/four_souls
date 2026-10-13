package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sucker (Basic Monster Card)
//
//	When this dies, it deals 1 damage to each player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sucker = engine.CardDef{
	Ref:     "sucker",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
