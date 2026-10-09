package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Xl Floor! (Good Event Card)
//
//	Expand monster slots by 1.
//	The active player may attack an additional time this turn.
var xlFloor = engine.CardDef{
	Ref:    "xl_floor",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Expand monster slots by 1. The active player may attack an additional time this turn.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandMonsters(1) }), engine.AddAttacks(1, engine.You)},
		},
	},
}
