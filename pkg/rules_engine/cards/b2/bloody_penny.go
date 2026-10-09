package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bloody Penny (Trinket Card)
//
//	Each time a player dies, before paying penalties, loot 1.
//	-Trinket- This loot becomes an item under your control when it resolves.
var bloodyPenny = engine.CardDef{
	Ref:     "bloody_penny",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player dies, before paying penalties, loot 1.",
			Trigger: engine.WhenAPlayerDies(),
			Effects: []engine.Effect{engine.Loot(1)},
		},
	},
}
