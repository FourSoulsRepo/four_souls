package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// I. The Magician (Wildcard Card)
//
//	Change the result of a dice roll to a number of your choosing.
var iTheMagician = engine.CardDef{
	Ref:    "i_the_magician",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Change the result of a dice roll to a number of your choosing.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				c.G.SetRoll(c.Targets[0].StackID, a[0]+1)
			}, engine.Question{
				Text:    "Change the roll to?",
				Options: func(*engine.Ctx, []int) []string { return []string{"1", "2", "3", "4", "5", "6"} },
			})},
		},
	},
}
