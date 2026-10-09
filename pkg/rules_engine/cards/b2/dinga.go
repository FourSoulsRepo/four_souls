package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dinga (Basic Monster Card)
//
//	When this dies on an attack roll of 6, double its rewards.
var dinga = engine.CardDef{
	Ref:     "dinga",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "When this dies on an attack roll of 6, double its rewards.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Prev == self && g.Object(self).KilledOn == 6
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Object(c.EventObject).DoubleRewards = true
			})},
		},
	},
}
