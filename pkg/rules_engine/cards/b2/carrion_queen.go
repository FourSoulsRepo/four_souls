package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Carrion Queen (Boss Card)
//
//	This takes no combat damage on attack rolls of 4 or 5.
var carrionQueen = engine.CardDef{
	Ref:     "carrion_queen",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if hit && (roll == 4 || roll == 5) {
			return 0 // no combat damage on attack rolls of 4 or 5
		}
		return n
	},
}
