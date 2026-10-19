package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Birthright (Active Treasure Card)
//
//	{Tap Effect}This becomes a copy of an eternal item you control. This loses eternal.
var birthright = engine.CardDef{
	Ref:    "birthright",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: This becomes a copy of an eternal item you control. This loses eternal.",
		Costs: []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetYourItem, func(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
			return c.Object != self && g.Eternal(c.Object)
		})},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			o := c.G.Object(c.Source)
			if o.Zone.Kind != engine.ZoneInPlay {
				return
			}
			o.CopyOf, o.CopyThisTurn, o.CopyTillYourTurn = c.G.CardOf(c.Targets[0].Object), false, false
			o.LostEternal = true
		})},
	}},
}
