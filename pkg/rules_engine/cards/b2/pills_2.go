package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Gain 4¢.
//	3-4: Gain 7¢.
//	5-6: Lose 4¢.
var pills2 = engine.CardDef{
	Ref:    "pills_2",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Roll- 1-2: Gain 4¢. 3-4: Gain 7¢. 5-6: Lose 4¢.",
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainCents(4)).
				Results(3, 4, engine.GainCents(7)).
				Results(5, 6, engine.LoseCents(4)))},
		},
	},
}
