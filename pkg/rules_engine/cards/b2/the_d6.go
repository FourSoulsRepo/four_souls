package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D6 (Eternal Treasure Card)
//
//	{Tap Effect}Choose a dice roll. Its controller rerolls it.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var theD6 = engine.CardDef{
	Ref:     "the_d6",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a dice roll. Its controller rerolls it.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.RerollRoll(0)},
		},
		rechargeAtEndOfTurn,
	},
}
