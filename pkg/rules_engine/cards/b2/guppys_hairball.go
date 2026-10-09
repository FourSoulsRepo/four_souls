package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Hairball (Trinket Card)
//
//	Each time you would take damage, roll-
//	6: Prevent 1 of that damage.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
//	-Trinket- This loot becomes an item under your control when it resolves.
var guppysHairball = engine.CardDef{
	Ref:     "guppys_hairball",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you would take damage, roll- 6: Prevent 1 of that damage.",
			Trigger: engine.WhenYouWouldTakeDamage(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(6, 6, engine.PreventDamage(1, engine.You)))},
		},
	},
}
