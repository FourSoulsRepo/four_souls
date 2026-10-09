package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// 3 Cents! (3 Cent Card)
//
//	Gain 3¢.
var threeCents = engine.CardDef{
	Ref:    "three_cents",
	Kind:   engine.LootCard,
	Copies: 15,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain 3¢.",
			Effects: []engine.Effect{engine.GainCents(3)},
		},
	},
}
