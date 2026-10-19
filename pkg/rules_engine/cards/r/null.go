package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Null (Basic Monster Card)
//
//	This takes no combat damage on every other attack roll made against it.
var null = engine.CardDef{
	Ref:     "null",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	CombatMod: func(g *engine.Game, _ engine.ObjectID, _ int, hit bool, n int) int {
		if hit && g.Attack.Rolls%2 == 0 {
			return 0 // every other attack roll: the 2nd, 4th, …
		}
		return n
	},
}
