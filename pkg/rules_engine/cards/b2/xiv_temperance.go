package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XIV. Temperance (Wildcard Card)
//
//	Choose one- Take 1 damage and gain 4¢. Take 2 damage and gain 8¢.
var xivTemperance = engine.CardDef{
	Ref:    "xiv_temperance",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Choose one- Take 1 damage and gain 4¢. Take 2 damage and gain 8¢.",
			Modes: []engine.Mode{
				{Text: "Take 1 damage and gain 4¢.", Effects: []engine.Effect{engine.DealDamage(1, engine.You), engine.GainCents(4)}},
				{Text: "Take 2 damage and gain 8¢.", Effects: []engine.Effect{engine.DealDamage(2, engine.You), engine.GainCents(8)}},
			},
		},
	},
}
