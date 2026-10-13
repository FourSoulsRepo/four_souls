package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Flaming Hopper (Basic Monster Card)
//
//	This takes no combat damage on attack rolls of 6.
var flamingHopper = engine.CardDef{
	Ref:     "flaming_hopper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if hit && roll == 6 {
			return 0 // "This takes no combat damage on attack rolls of 6."
		}
		return n
	},
}
