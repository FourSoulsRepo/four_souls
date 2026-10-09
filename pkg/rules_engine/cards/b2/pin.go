package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pin (Boss Card)
//
//	This takes no combat damage on attack rolls of 6.
var pin = engine.CardDef{
	Ref:     "pin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if hit && roll == 6 {
			return 0
		}
		return n
	},
}
