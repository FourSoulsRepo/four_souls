package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Wishbone (Trinket Card)
//
//	When this is destroyed, gain +1 treasure.
//	-Trinket- This loot becomes an item under your control when it resolves.
var wishbone = engine.CardDef{
	Ref:     "wishbone",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this is destroyed, gain +1 treasure.",
		Trigger: engine.WhenThisIsDestroyed(),
		Effects: []engine.Effect{engine.GainTreasure(1)},
	}},
}
