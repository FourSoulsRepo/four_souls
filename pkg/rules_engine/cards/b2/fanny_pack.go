package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fanny Pack (Passive Treasure Card)
//
//	Each time you take damage, loot 1.
var fannyPack = engine.CardDef{
	Ref:    "fanny_pack",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, loot 1.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.Loot(1)},
		},
	},
}
