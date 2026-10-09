package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dople (Basic Monster Card)
//
//	Damage dealt to this is also dealt to the player to the active player's right.
var dople = engine.CardDef{
	Ref:     "dople",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Damage dealt to this is also dealt to the player to the active player's right.",
			Trigger: whenThisTakesDamage(false),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				p := c.G.PlayerBefore(c.G.Turn.Active)
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, c.EventAmount, engine.NoPlayer, c.Source)
			})},
		},
	},
}
