package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cod Worm (Basic Monster Card)
//
// No card text: stats and rewards only.
var codWorm = engine.CardDef{
	Ref:     "cod_worm",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     0,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
