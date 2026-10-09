package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hanger (Basic Monster Card)
//
//	When this dies, expand shop slots by 1.
var hanger = engine.CardDef{
	Ref:     "hanger",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, expand shop slots by 1.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.ExpandShop(1)
			})},
		},
	},
}
