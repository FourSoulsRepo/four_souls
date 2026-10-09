package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Clotty (Basic Monster Card)
//
// No card text: stats and rewards only.
var clotty = engine.CardDef{
	Ref:     "clotty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
}
