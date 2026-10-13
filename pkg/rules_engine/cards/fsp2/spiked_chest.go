package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Spiked Chest (Good Event Card)
//
//	Roll-
//	1-2: Take 1 damage.
//	3-4: Take 1 damage. Loot 1.
//	5-6: Take 1 damage. Gain +1 treasure.
var spikedChest = engine.CardDef{
	Ref:    "spiked_chest",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Roll- 1-2: Take 1 damage. 3-4: Take 1 damage. Loot 1. 5-6: Take 1 damage. Gain +1 treasure.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.DealDamage(1, engine.You)).
			Results(3, 4, engine.DealDamage(1, engine.You), engine.Loot(1)).
			Results(5, 6, engine.DealDamage(1, engine.You), engine.GainTreasure(1)))},
	}},
}
