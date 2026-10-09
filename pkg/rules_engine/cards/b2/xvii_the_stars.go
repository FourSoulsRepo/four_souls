package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XVII. The Stars (Wildcard Card)
//
//	Gain +1 treasure.
var xviiTheStars = engine.CardDef{
	Ref:    "xvii_the_stars",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Gain +1 treasure.",
			Effects: []engine.Effect{engine.GainTreasure(1)},
		},
	},
}
