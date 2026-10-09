package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gold Chest (Good Event Card)
//
//	Roll-
//	1-2: Gain +1 Treasure.
//	3-4: Loot 1.
//	5-6: Loot 2.
var goldChest2 = engine.CardDef{
	Ref:    "gold_chest_2",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-2: Gain +1 treasure. 3-4: Loot 1. 5-6: Loot 2.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainTreasure(1)).
				Results(3, 4, engine.Loot(1)).
				Results(5, 6, engine.Loot(2)))},
		},
	},
}
