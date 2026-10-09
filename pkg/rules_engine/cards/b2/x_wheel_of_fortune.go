package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// X. Wheel Of Fortune (Wildcard Card)
//
//	Roll-
//	1: Gain 1¢.
//	2: Take 2 damage.
//	3. Loot 3.
//	4. Lose 4¢.
//	5: Gain 5¢.
//	6: Gain +1 treasure.
var xWheelOfFortune = engine.CardDef{
	Ref:    "x_wheel_of_fortune",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Roll- 1: Gain 1¢. 2: Take 2 damage. 3: Loot 3. 4: Lose 4¢. 5: Gain 5¢. 6: Gain +1 treasure.",
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 1, engine.GainCents(1)).
				Results(2, 2, engine.DealDamage(2, engine.You)).
				Results(3, 3, engine.Loot(3)).
				Results(4, 4, engine.LoseCents(4)).
				Results(5, 5, engine.GainCents(5)).
				Results(6, 6, engine.GainTreasure(1)))},
		},
	},
}
