package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Tumor (Basic Monster Card)
//
//	Each time this takes damage, it deals 1 damage to each non-active player.
var tumor = engine.CardDef{
	Ref:     "tumor",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, it deals 1 damage to each non-active player.",
		Trigger: whenThisTakesDamage(false),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for _, p := range livingPlayers(c.G, c.G.Turn.Active) {
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			}
		})},
	}},
}
