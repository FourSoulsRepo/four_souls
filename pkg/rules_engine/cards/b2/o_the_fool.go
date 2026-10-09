package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// O. The Fool (Wildcard Card)
//
//	End the turn. Cancel everything that hasn't resolved.
var oTheFool = engine.CardDef{
	Ref:    "o_the_fool",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "End the turn. Cancel everything that hasn't resolved.",
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.EndTurnNow() })},
		},
	},
}
