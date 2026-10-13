package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Nerve Ending (Basic Monster Card)
//
// No card text: stats and rewards only.
var nerveEnding = engine.CardDef{
	Ref:     "nerve_ending",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
