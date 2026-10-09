package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pale Fatty (Basic Monster Card)
//
// No card text: stats and rewards only.
var paleFatty = engine.CardDef{
	Ref:     "pale_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
