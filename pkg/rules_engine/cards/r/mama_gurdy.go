package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mama Gurdy (Boss Card)
//
// No text: a soul and "Roll- loot x" rewards.
var mamaGurdy = engine.CardDef{
	Ref:     "mama_gurdy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      7,
	DC:      3,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
}
