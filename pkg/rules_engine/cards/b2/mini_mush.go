package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mini Mush (Active Treasure Card)
//
//	{Tap Effect}Subtract up to 2 from a roll.
var miniMush = engine.CardDef{
	Ref:    "mini_mush",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Subtract up to 2 from a roll.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{changeRollBy(-1, -2)},
		},
	},
}
