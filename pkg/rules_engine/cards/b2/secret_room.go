package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Secret Room! (Good Event Card)
//
//	Roll-
//	1: Take 3 damage.
//	2-3: Discard 2 loot cards.
//	4-5: Gain 7¢.
//	6: Gain +1 Treasure.
var secretRoom = engine.CardDef{
	Ref:    "secret_room",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Roll- 1: Take 3 damage. 2-3: Discard 2 loot cards. 4-5: Gain 7¢. 6: Gain +1 treasure.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 1, engine.DealDamage(3, engine.You)).
				Results(2, 3, discardOne, discardOne).
				Results(4, 5, engine.GainCents(7)).
				Results(6, 6, engine.GainTreasure(1)))},
		},
	},
}
