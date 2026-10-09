package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Modeling Clay (Active Treasure Card)
//
//	{Tap Effect}Choose a non-eternal item. This becomes a copy of that item.
//	(This change is indefinite.)
var modelingClay = engine.CardDef{
	Ref:    "modeling_clay",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a non-eternal item. This becomes a copy of that item.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, otherNonEternalItem)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Object(c.Source).CopyOf = c.G.CardOf(c.Targets[0].Object)
			})},
		},
	},
}

// otherNonEternalItem keeps non-eternal items other than this one.
func otherNonEternalItem(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
	return c.Object != self && !g.Eternal(c.Object)
}
