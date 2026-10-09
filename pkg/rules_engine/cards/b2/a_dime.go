package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Dime!! (Nickel Card)
//
//	Gain 10¢.
var aDime = engine.CardDef{
	Ref:    "a_dime",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain 10¢.",
			Effects: []engine.Effect{engine.GainCents(10)},
		},
	},
}
