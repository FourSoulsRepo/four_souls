package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Conquest (Boss Card)
//
//	When this dies, the active player must make an additional attack.
var conquest = engine.CardDef{
	Ref:     "conquest",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
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
