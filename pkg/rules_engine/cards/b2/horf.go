package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Horf (Basic Monster Card)
//
//	Combat damage this deals is increased by 1 on attack rolls of 2.
var horf = engine.CardDef{
	Ref:     "horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if !hit && roll == 2 {
			return n + 1
		}
		return n
	},
}
