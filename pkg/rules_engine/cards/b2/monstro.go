package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Monstro (Boss Card)
//
// No card text: stats and rewards only.
var monstro = engine.CardDef{
	Ref:     "monstro",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
