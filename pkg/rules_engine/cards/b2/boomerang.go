package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Boomerang (Active Treasure Card)
//
//	{Tap Effect}Choose another player. Steal a loot card from them at random.
var boomerang = engine.CardDef{
	Ref:    "boomerang",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose another player. Steal a loot card from them at random.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					from := c.Targets[0].Player
					c.G.GiveHandCard(from, c.Controller, c.G.Players[from].Hand[a[0]])
				}
			}, engine.Question{Random: true, Text: "Which card?", Options: func(c *engine.Ctx, _ []int) []string {
				return make([]string, len(c.G.Players[c.Targets[0].Player].Hand)) // hidden: labels stay empty
			}})},
		},
	},
}
