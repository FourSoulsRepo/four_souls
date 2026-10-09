package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fly (Basic Monster Card)
//
// No card text: stats and rewards only.
var fly = engine.CardDef{
	Ref:     "fly",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
}
