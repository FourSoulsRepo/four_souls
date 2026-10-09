package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Godhead (Active Treasure Card)
//
//	{Tap Effect}Change the result of a dice roll to a 1 or 6.
var godhead = engine.CardDef{
	Ref:    "godhead",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Change the result of a dice roll to a 1 or 6.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				c.G.SetRoll(c.Targets[0].StackID, []int{1, 6}[a[0]])
			}, engine.Question{Text: "Change it to?", Options: func(*engine.Ctx, []int) []string { return []string{"1", "6"} }})},
		},
	},
}
