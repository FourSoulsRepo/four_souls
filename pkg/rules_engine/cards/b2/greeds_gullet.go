package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed’s Gullet (Passive Treasure Card)
//
//	Each time you die, before paying penalties, gain 8¢.
var greedsGullet = engine.CardDef{
	Ref:    "greeds_gullet",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you die, before paying penalties, gain 8¢.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{engine.GainCents(8)},
		},
	},
}
