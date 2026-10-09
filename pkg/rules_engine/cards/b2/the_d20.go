package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D20 (Active Treasure Card)
//
//	{Tap Effect}Reroll an item.
//	(Destroy that item and replace it with the top card of the treasure deck.)
var theD20 = engine.CardDef{
	Ref:    "the_d20",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Reroll an item.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetItem)},
			Effects: []engine.Effect{engine.RerollTarget(0)},
		},
	},
}
