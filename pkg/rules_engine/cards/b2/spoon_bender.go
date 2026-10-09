package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Spoon Bender (Active Treasure Card)
//
//	{Tap Effect}Add 1 to a roll.
var spoonBender = engine.CardDef{
	Ref:    "spoon_bender",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Add 1 to a roll.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.ModifyRoll(1, 0)},
		},
	},
}
