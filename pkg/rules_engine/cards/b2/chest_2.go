package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chest (Good Event Card)
//
//	Roll-
//	1-2: Loot 1.
//	3-4: Loot 2.
//	5-6: Loot 3.
var chest2 = engine.CardDef{
	Ref:    "chest_2",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-2: Loot 1. 3-4: Loot 2. 5-6: Loot 3.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.Loot(1)).
				Results(3, 4, engine.Loot(2)).
				Results(5, 6, engine.Loot(3)))},
		},
	},
}
