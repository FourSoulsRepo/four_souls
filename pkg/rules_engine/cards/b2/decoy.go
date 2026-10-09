package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Decoy (Active Treasure Card)
//
//	{Tap Effect}Swap this with a non-eternal item another player controls.
var decoy = engine.CardDef{
	Ref:    "decoy",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Swap this with a non-eternal item another player controls.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, othersNonEternal)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				them := c.Targets[0].Player
				c.G.GainControl(c.Controller, c.Targets[0].Object)
				c.G.GainControl(them, c.Source)
			})},
		},
	},
}

// othersNonEternal keeps non-eternal items of other players.
func othersNonEternal(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
	return c.Player != g.Object(self).Controller && !g.Eternal(c.Object)
}
