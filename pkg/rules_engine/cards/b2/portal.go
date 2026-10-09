package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Portal (Basic Monster Card)
//
//	When this dies, the active player must make an additional attack.
var portal = engine.CardDef{
	Ref:     "portal",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player must make an additional attack.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.ForceAttacks(1, false)
			})},
		},
	},
}
