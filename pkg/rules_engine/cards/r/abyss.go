package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Abyss (Eternal Treasure Card)
//
//	{Tap Effect}Destroy an item you control. If you do, put a counter on this.
//	You have +1{ATK} for every 2 counters on this.
//	-Eternal- This can't be destroyed or put into discard.
var abyss = engine.CardDef{
	Ref:     "abyss",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Destroy an item you control. If you do, put a counter on this.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetYourItem, engine.NotThis)},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if c.G.DestroyObject(c.Controller, c.Targets[0].Object) {
				c.G.AddCountersTo(c.Source, 1)
			}
		})},
	}},
	Statics: []engine.Static{{
		Stat: engine.StatPlayerATK,
		Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
			return g.Object(self).Controller == p
		},
		Value: func(g *engine.Game, self engine.ObjectID) int { return g.Object(self).CountersOf("") / 2 },
	}},
}
