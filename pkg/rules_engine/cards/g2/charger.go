package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Charger (Basic Monster Card)
//
// No card text: stats and rewards only.
var charger = engine.CardDef{
	Ref:     "charger",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
