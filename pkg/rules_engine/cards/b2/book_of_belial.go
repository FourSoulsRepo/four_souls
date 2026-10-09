package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Book of Belial (Eternal Treasure Card)
//
//	Add or subtract 1 from a roll.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// The printed card has the ↷ icon (the site text lost it): it is used
// like the other recharging starting items.
var bookOfBelial = engine.CardDef{
	Ref:     "book_of_belial",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Add or subtract 1 from a roll.",
			Costs: []engine.Cost{engine.Tap()},
			Modes: []engine.Mode{
				{
					Text:    "Add 1 to a roll.",
					Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
					Effects: []engine.Effect{engine.ModifyRoll(1, 0)},
				},
				{
					Text:    "Subtract 1 from a roll.",
					Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
					Effects: []engine.Effect{engine.ModifyRoll(-1, 0)},
				},
			},
		},
		rechargeAtEndOfTurn,
	},
}
