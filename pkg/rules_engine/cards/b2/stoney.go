package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Stoney (Basic Monster Card)
//
//	Monsters have +1{DC}.
//	This can't be attacked.
//	When another monster dies, this dies.
var stoney = engine.CardDef{
	Ref:          "stoney",
	Kind:         engine.MonsterCard,
	Copies:       1,
	HP:           3,
	DC:           0,
	ATK:          0,
	Rewards:      []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Unattackable: true,
	Statics:      []engine.Static{engine.MonstersHave(engine.StatMonsterDC, 1)},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "When another monster dies, this dies.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Prev != self && g.Kind(e.Object) == engine.MonsterCard && g.Object(self).Zone.Kind == engine.ZoneInPlay
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.KillObject(c.Source)
			})},
		},
	},
}
