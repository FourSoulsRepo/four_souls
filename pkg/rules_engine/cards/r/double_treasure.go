package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Double Treasure! (Good Event Card)
//
//	Choose another player. You and that player each gain +1 treasure.
var doubleTreasure = engine.CardDef{
	Ref:    "double_treasure",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Choose another player. You and that player each gain +1 treasure.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			ps := othersOfActive(c.G)
			if a[0] < 0 || a[0] >= len(ps) {
				return
			}
			c.Do(engine.GainTreasure(1))
			them := *c
			them.Controller = ps[a[0]]
			them.Do(engine.GainTreasure(1))
		}, engine.Question{Text: "Choose another player.", Options: func(c *engine.Ctx, _ []int) []string {
			return playerLabels(c.G, othersOfActive(c.G))
		}})},
	}},
}
