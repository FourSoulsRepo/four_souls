package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dressing Table (Good Event Card)
//
//	Reroll any number of items you control.
var dressingTable = engine.CardDef{
	Ref:    "dressing_table",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Reroll any number of items you control.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{rerollAnyNumber(8)},
	}},
}
