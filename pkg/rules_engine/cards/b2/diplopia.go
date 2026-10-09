package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Diplopia (Active Treasure Card)
//
//	{Tap Effect}Choose a non-eternal passive item. This becomes a copy of that item till end of turn.
var diplopia = engine.CardDef{
	Ref:    "diplopia",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a non-eternal passive item. This becomes a copy of that item till end of turn.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, passiveNonEternal)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				o := c.G.Object(c.Source)
				o.CopyOf, o.CopyThisTurn = c.G.CardOf(c.Targets[0].Object), true
			})},
		},
	},
}

// passiveNonEternal keeps other non-eternal items without activated
// abilities.
func passiveNonEternal(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
	if c.Object == self || g.Eternal(c.Object) {
		return false
	}
	for _, ref := range g.AbilitiesOf(c.Object) {
		if g.Ability(ref).Kind == engine.Activated {
			return false
		}
	}
	return true
}
