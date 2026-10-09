package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Troll Bombs (Bad Event Card)
//
//	Take 2 damage!
var trollBombs = engine.CardDef{
	Ref:    "troll_bombs",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Take 2 damage!",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.DealDamage(2, engine.You)},
		},
	},
}
