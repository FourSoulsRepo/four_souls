package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Widow (Boss Card)
//
// No card text: stats and rewards only.
var widow = engine.CardDef{
	Ref:     "widow",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
