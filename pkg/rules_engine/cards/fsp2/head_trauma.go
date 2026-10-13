package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Head Trauma (Bad Event Card)
//
//	Discard 2 loot cards.
var headTrauma = engine.CardDef{
	Ref:    "head_trauma",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Discard 2 loot cards.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{discardOne, discardOne},
	}},
}
