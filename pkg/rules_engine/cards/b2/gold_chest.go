package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gold Chest (Good Event Card)
//
//	Roll-
//	1-2: Gain +1 Treasure.
//	3-4: Gain 5¢.
//	5-6: Gain 7¢.
var goldChest = engine.CardDef{
	Ref:    "gold_chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-2: Gain +1 treasure. 3-4: Gain 5¢. 5-6: Gain 7¢.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainTreasure(1)).
				Results(3, 4, engine.GainCents(5)).
				Results(5, 6, engine.GainCents(7)))},
		},
	},
}
