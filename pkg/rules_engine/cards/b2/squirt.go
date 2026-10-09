package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Squirt (Basic Monster Card)
//
// No card text: stats and rewards only.
var squirt = engine.CardDef{
	Ref:     "squirt",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
