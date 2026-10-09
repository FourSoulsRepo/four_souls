package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Bony (Basic Monster Card)
//
//	When this dies, it deals 1 damage to the player who killed it.
var blackBony = engine.CardDef{
	Ref:     "black_bony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, it deals 1 damage to the player who killed it.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if k := c.G.Object(c.Source).Killer(); k != engine.NoPlayer {
					c.G.DealDamageTo(engine.Target{Player: k, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
				}
			})},
		},
	},
}
