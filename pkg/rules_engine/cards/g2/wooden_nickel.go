package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Wooden Nickel (Eternal Treasure Card)
//
//	{Tap Effect}Choose a player, then roll-
//	That player gains ¢ equal to the result.
//	-Eternal- This can't be destroyed or put into discard.
var woodenNickel = engine.CardDef{
	Ref:     "wooden_nickel",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a player, then roll- That player gains ¢ equal to the result.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.Roll(centsEqualToRoll())},
		},
	},
}
