package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fat Bat (Basic Monster Card)
//
// No card text: stats and rewards only.
var fatBat = engine.CardDef{
	Ref:     "fat_bat",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
