package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Nickel! (Nickel Card)
//
//	Gain 5¢
var aNickel = engine.CardDef{
	Ref:    "a_nickel",
	Kind:   engine.LootCard,
	Copies: 5,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain 5¢.",
			Effects: []engine.Effect{engine.GainCents(5)},
		},
	},
}
