package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mulligan (Basic Monster Card)
//
//	When this dies, expand monster slots by 1.
var mulligan = engine.CardDef{
	Ref:     "mulligan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
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
