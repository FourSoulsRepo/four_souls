package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Leaper (Basic Monster Card)
//
//	Combat damage this deals is doubled on attack rolls of 1.
var leaper = engine.CardDef{
	Ref:     "leaper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if !hit && roll == 1 {
			return n * 2
		}
		return n
	},
}
