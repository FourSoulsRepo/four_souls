package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Book Of Sin (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1-2: Gain 1¢.
//	3-4: Loot 1.
//	5-6: Gain +1{HP} till end of turn.
var bookOfSin = engine.CardDef{
	Ref:    "book_of_sin",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Roll- 1-2: Gain 1¢. 3-4: Loot 1. 5-6: Gain +1 HP till end of turn.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainCents(1)).
				Results(3, 4, engine.Loot(1)).
				Results(5, 6, engine.GainHPThisTurn(1, engine.You)))},
		},
	},
}
