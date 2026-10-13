package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// I Am Error! (Good Event Card)
//
//	Roll-
//	1-2: You must make an additional attack on the monster deck.
//	3-4: Loot 3.
//	5-6: Gain +1 Treasure.
var iAmError = engine.CardDef{
	Ref:    "i_am_error",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Roll- 1-2: You must make an additional attack on the monster deck. 3-4: Loot 3. 5-6: Gain +1 treasure.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.EffectFunc(func(c *engine.Ctx) { c.G.ForceAttacks(1, true) })).
			Results(3, 4, engine.Loot(3)).
			Results(5, 6, engine.GainTreasure(1)))},
	}},
}
