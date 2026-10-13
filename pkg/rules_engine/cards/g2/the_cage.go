package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Cage (Boss Card)
//
// No card text: stats and rewards only.
var theCage = engine.CardDef{
	Ref:     "the_cage",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
