package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chest (Good Event Card)
//
//	Roll-
//	1-2: Gain 1¢.
//	3-4: Gain 3¢.
//	5-6: Gain 6¢.
var chest = engine.CardDef{
	Ref:    "chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-2: Gain 1¢. 3-4: Gain 3¢. 5-6: Gain 6¢.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainCents(1)).
				Results(3, 4, engine.GainCents(3)).
				Results(5, 6, engine.GainCents(6)))},
		},
	},
}
