package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// 2 Cents! (2 Cent Card)
//
//	Gain 2¢.
var twoCents = engine.CardDef{
	Ref:    "two_cents",
	Kind:   engine.LootCard,
	Copies: 12,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain 2¢.",
			Effects: []engine.Effect{engine.GainCents(2)},
		},
	},
}
