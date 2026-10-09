package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Suicide King (Passive Treasure Card)
//
//	Each time you die, before paying penalties, loot 3.
var suicideKing = engine.CardDef{
	Ref:    "suicide_king",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you die, before paying penalties, loot 3.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{engine.Loot(3)},
		},
	},
}
