package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom! (Epic Boss Card)
//
//	Combat damage this deals is doubled on attack rolls of 1.
//	When this dies, expand monsters slots by 1.
var mom = engine.CardDef{
	Ref:     "mom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    2,
	HP:      5,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if !hit && roll == 1 {
			return n * 2
		}
		return n
	},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, expand monster slots by 1.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.ExpandMonsters(1)
			})},
		},
	},
}
