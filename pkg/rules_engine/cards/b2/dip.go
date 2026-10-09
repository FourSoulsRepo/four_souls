package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dip (Basic Monster Card)
//
// No card text: stats and rewards only.
var dip = engine.CardDef{
	Ref:     "dip",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
}
