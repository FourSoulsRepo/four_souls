package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Evil Twin (Basic Monster Card)
//
//	Damage dealt to this is also dealt to the player to the active player's left.
var evilTwin = engine.CardDef{
	Ref:     "evil_twin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Damage dealt to this is also dealt to the player to the active player's left.",
			Trigger: whenThisTakesDamage(false),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				p := c.G.PlayerAfter(c.G.Turn.Active)
				c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, c.EventAmount, engine.NoPlayer, c.Source)
			})},
		},
	},
}
