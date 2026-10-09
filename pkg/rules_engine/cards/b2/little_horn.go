package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Little Horn (Boss Card)
//
// No card text: stats and rewards only.
var littleHorn = engine.CardDef{
	Ref:     "little_horn",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      6,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
