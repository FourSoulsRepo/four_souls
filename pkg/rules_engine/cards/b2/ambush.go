package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ambush! (Bad Event Card)
//
//	The active player must attack the monster deck 2 times this turn.
var ambush = engine.CardDef{
	Ref:    "ambush",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "The active player must attack the monster deck 2 times this turn.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ForceAttacks(2, true) })},
		},
	},
}
