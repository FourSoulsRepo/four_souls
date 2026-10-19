package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Gutted Fatty (Basic Monster Card)
//
//	When this dies expand monster slots by 1. The active player must make an additional attack.
var guttedFatty = engine.CardDef{
	Ref:     "gutted_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies expand monster slots by 1. The active player must make an additional attack.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.ExpandMonsters(1)
			c.G.ForceAttacks(1, false)
		})},
	}},
}
