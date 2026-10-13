package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Penny! (1 Cent Card)
//
//	Gain 1¢.
var aPenny = engine.CardDef{
	Ref:    "a_penny",
	Kind:   engine.LootCard,
	Copies: 6,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain 1¢.",
			Effects: []engine.Effect{engine.GainCents(1)},
		},
	},
}
