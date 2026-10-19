package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pitfall! (Bad Event Card)
//
//	End the turn. Cancel everything that hasn't resolved.
var pitfall = engine.CardDef{
	Ref:    "pitfall",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "End the turn. Cancel everything that hasn't resolved.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.EndTurnNow() })},
	}},
}
