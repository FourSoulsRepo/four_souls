package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hush (Epic Boss Card)
//
// No card text: stats and rewards only.
var hush = engine.CardDef{
	Ref:     "hush",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      8,
	DC:      3,
	ATK:     2,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
