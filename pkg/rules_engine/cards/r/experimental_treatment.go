package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Experimental Treatment (Soul Treasure Card)
//
//	{Paid Effect}Destroy this: roll-
//	1: You die.
//	2: Loot 3.
//	3: Gain 9¢.
//	4: Gain +1 treasure.
//	5: Gain +2 treasure.
//	6: This becomes a soul. Gain it.
var experimentalTreatment = engine.CardDef{
	Ref:    "experimental_treatment",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Destroy this: roll- 1: You die. 2: Loot 3. 3: Gain 9¢. 4: Gain +1 treasure. 5: Gain +2 treasure. 6: This becomes a soul. Gain it.",
		Costs: []engine.Cost{engine.DestroySelf()},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 1, killYou).
			Results(2, 2, engine.Loot(3)).
			Results(3, 3, engine.GainCents(9)).
			Results(4, 4, engine.GainTreasure(1)).
			Results(5, 5, engine.GainTreasure(2)).
			// As a soul it keeps this ability (R-CARD-38, S-YUG).
			Results(6, 6, engine.EffectFunc(func(c *engine.Ctx) {
				c.G.SoulFromDiscard(engine.TreasureDeck, "experimental_treatment", c.Controller)
			})))},
	}},
}
