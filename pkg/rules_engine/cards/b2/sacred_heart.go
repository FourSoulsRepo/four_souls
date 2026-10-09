package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sacred Heart (Passive Treasure Card)
//
//	When you would roll a 1, you may change the result to a 6.
var sacredHeart = engine.CardDef{
	Ref:    "sacred_heart",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "When you would roll a 1, you may change the result to a 6.",
			Trigger: engine.Trigger{On: engine.EvRollWouldResolve, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Amount == 1 && e.Player == g.Object(self).Controller
			}},
			Effects: []engine.Effect{may("Change the 1 to a 6?", engine.EffectFunc(func(c *engine.Ctx) { c.G.SetRoll(c.EventStack, 6) }))},
		},
	},
}
