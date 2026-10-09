package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hopper (Basic Monster Card)
//
//	This takes no combat damage on attack rolls of 6.
var hopper = engine.CardDef{
	Ref:     "hopper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if hit && roll == 6 {
			return 0
		}
		return n
	},
}
