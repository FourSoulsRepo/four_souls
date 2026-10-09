package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Trite (Basic Monster Card)
//
// No card text: stats and rewards only.
var trite = engine.CardDef{
	Ref:     "trite",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
