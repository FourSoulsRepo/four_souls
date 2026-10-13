package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Angel Room (Good Event Card)
//
//	Roll-
//	1: Gain +2 treasure.
//	2-3: Gain +1 treasure.
//	4-6: Loot 2.
var angelRoom = engine.CardDef{
	Ref:    "angel_room",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Roll- 1: Gain +2 treasure. 2-3: Gain +1 treasure. 4-6: Loot 2.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 1, engine.GainTreasure(2)).
			Results(2, 3, engine.GainTreasure(1)).
			Results(4, 6, engine.Loot(2)))},
	}},
}
