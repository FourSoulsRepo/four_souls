package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mystery Sack (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1-2: Loot 1.
//	3-4: Gain 4¢.
var mysterySack = engine.CardDef{
	Ref:    "mystery_sack",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Roll- 1-2: Loot 1. 3-4: Gain 4¢.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 2, engine.Loot(1)).Results(3, 4, engine.GainCents(4)))},
		},
	},
}
