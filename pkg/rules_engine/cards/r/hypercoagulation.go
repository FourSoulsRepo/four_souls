package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hypercoagulation (Eternal Treasure Card)
//
//	At the start of your turn put a counter on this.
//	You have +1{HP} for each counter on this.
//	At the end of your turn, if this has 3+ counters, remove all of them and loot 3.
//	-Eternal- This can't be destroyed or put into discard.
var hypercoagulation = engine.CardDef{
	Ref:     "hypercoagulation",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn put a counter on this.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind: engine.Triggered,
			Text: "At the end of your turn, if this has 3+ counters, remove all of them and loot 3.",
			Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				o := g.Object(self)
				return e.Player == o.Controller && o.CountersOf("") >= 3
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if n := c.G.Object(c.Source).CountersOf(""); n >= 3 {
					c.G.AddCountersTo(c.Source, -n)
					c.Do(engine.Loot(3))
				}
			})},
		},
	},
	Statics: []engine.Static{{
		Stat: engine.StatPlayerHP,
		Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
			return g.Object(self).Controller == p
		},
		Value: func(g *engine.Game, self engine.ObjectID) int { return g.Object(self).CountersOf("") },
	}},
}
