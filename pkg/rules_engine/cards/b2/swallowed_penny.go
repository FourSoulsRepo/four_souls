package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Swallowed Penny (Trinket Card)
//
//	Each time you take damage, gain 1¢.
//	-Trinket- This loot becomes an item under your control when it resolves.
var swallowedPenny = engine.CardDef{
	Ref:     "swallowed_penny",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, gain 1¢.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.GainCents(1)},
		},
	},
}
