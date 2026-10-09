package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fatty (Basic Monster Card)
//
// No card text: stats and rewards only.
var fatty = engine.CardDef{
	Ref:     "fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
