package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper Head (Basic Monster Card)
//
//	Each time this deals combat damage to a player, they lose 2¢.
var keeperHead = engine.CardDef{
	Ref:     "keeper_head",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals combat damage to a player, they lose 2¢.",
			Trigger: whenThisDealsDamage(true),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.LoseCentsNow(c.EventPlayer, 2)
			})},
		},
	},
}
