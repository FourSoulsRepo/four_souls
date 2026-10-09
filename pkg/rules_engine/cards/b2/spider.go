package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Spider (Basic Monster Card)
//
// No card text: stats and rewards only.
var spider = engine.CardDef{
	Ref:     "spider",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
