package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Chest (Good Event Card)
//
//	Roll-
//	1-2: Loot 1.
//	3-4: Gain 3¢.
//	5-6: Take 2 damage.
var darkChest = engine.CardDef{
	Ref:    "dark_chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-2: Loot 1. 3-4: Gain 3¢. 5-6: Take 2 damage.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.Loot(1)).
				Results(3, 4, engine.GainCents(3)).
				Results(5, 6, engine.DealDamage(2, engine.You)))},
		},
	},
}
