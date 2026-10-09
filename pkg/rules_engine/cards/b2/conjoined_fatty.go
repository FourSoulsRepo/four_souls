package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Conjoined Fatty (Basic Monster Card)
//
// No card text: stats and rewards only.
var conjoinedFatty = engine.CardDef{
	Ref:     "conjoined_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
