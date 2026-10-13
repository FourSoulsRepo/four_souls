package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Roundy (Basic Monster Card)
//
//	Each time this deals combat damage, it deals 1 damage to each other monster.
var roundy = engine.CardDef{
	Ref:     "roundy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this deals combat damage, it deals 1 damage to each other monster.",
		Trigger: whenThisDealsDamage(true),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for _, id := range monstersInPlay(c.G) {
				if id != c.Source {
					c.G.DealDamageTo(engine.Target{Object: id}, 1, engine.NoPlayer, c.Source)
				}
			}
		})},
	}},
}
