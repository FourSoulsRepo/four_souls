package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Leech (Basic Monster Card)
//
// No card text: stats and rewards only.
var leech = engine.CardDef{
	Ref:     "leech",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
