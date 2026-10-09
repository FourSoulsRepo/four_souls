package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Smelter (Paid Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Gain 3¢.
var smelter = engine.CardDef{
	Ref:    "smelter",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Discard a loot card: Gain 3¢.",
			Costs:   []engine.Cost{engine.DiscardChosen(0)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
			Effects: []engine.Effect{engine.GainCents(3)},
		},
	},
}
