package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Fatty (Boss Card)
//
//	Each time this deals combat damage, it heals 1{HP}.
var megaFatty = engine.CardDef{
	Ref:     "mega_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals combat damage, it heals 1 HP.",
			Trigger: whenThisDealsDamage(true),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.HealObject(c.Source, 1)
			})},
		},
	},
}
