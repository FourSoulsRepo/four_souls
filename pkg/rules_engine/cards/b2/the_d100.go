package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D100 (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1: Loot 1.
//	2: Loot 2.
//	3: Gain 3¢.
//	4: Gain 4¢.
//	5: Gain +1{HP} till end of turn.
//	6: Gain +1{ATK} till end of turn.
var theD100 = engine.CardDef{
	Ref:    "the_d100",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Roll- 1: Loot 1. 2: Loot 2. 3: Gain 3¢. 4: Gain 4¢. 5: Gain +1 HP till end of turn. 6: Gain +1 ATK till end of turn.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 1, engine.Loot(1)).
				Results(2, 2, engine.Loot(2)).
				Results(3, 3, engine.GainCents(3)).
				Results(4, 4, engine.GainCents(4)).
				Results(5, 5, engine.GainHPThisTurn(1, engine.You)).
				Results(6, 6, engine.GainATKThisTurn(1, engine.You)))},
		},
	},
}
