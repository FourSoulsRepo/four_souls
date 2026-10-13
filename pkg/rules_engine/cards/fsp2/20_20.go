package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// 20/20 (Active Treasure Card)
//
//	{Tap Effect}Add up to 2 to an attack roll.
var card2020 = engine.CardDef{
	Ref:    "20_20",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Add up to 2 to an attack roll.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetDiceRoll, isAttackRoll)},
		Effects: []engine.Effect{changeRollBy(1, 2)},
	}},
}
