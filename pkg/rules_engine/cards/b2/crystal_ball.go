package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Crystal Ball (Active Treasure Card)
//
//	{Tap Effect}Before a dice is rolled, choose a number. If the next roll is that number, loot 3.
var crystalBall = engine.CardDef{
	Ref:    "crystal_ball",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Dice: func(g *engine.Game, self engine.ObjectID) []int {
		if n := g.Object(self).Chosen; n > 0 {
			return []int{n}
		}
		return nil
	},
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Before a dice is rolled, choose a number. If the next roll is that number, loot 3.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				c.G.Object(c.Source).Chosen = a[0] + 1
			}, engine.Question{Text: "Which number?", Options: func(*engine.Ctx, []int) []string { return []string{"1", "2", "3", "4", "5", "6"} }})},
		},
		{
			Kind: engine.Triggered,
			Text: "If the next roll is that number, loot 3.",
			Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, _ engine.Event) bool {
				return g.Object(self).Chosen > 0
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				o := c.G.Object(c.Source)
				guess := o.Chosen
				o.Chosen = 0
				if guess == c.EventAmount {
					c.Do(engine.Loot(3))
				}
			})},
		},
	},
}
