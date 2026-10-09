package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// 4 Cents! (4 Cent Card)
//
//	Gain 4¢.
var fourCents = engine.CardDef{
	Ref:    "four_cents",
	Kind:   engine.LootCard,
	Copies: 9,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain 4¢.",
			Effects: []engine.Effect{engine.GainCents(4)},
		},
	},
}
