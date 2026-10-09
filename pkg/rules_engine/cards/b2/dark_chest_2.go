package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Chest (Good Event Card)
//
//	Roll-
//	1-2: Gain 1¢.
//	3-4: Loot 2.
//	5-6: Take 2 damage.
var darkChest2 = engine.CardDef{
	Ref:    "dark_chest_2",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1-2: Gain 1¢. 3-4: Loot 2. 5-6: Take 2 damage.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainCents(1)).
				Results(3, 4, engine.Loot(2)).
				Results(5, 6, engine.DealDamage(2, engine.You)))},
		},
	},
}
