package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Shovel (One-Use Treasure Card)
//
//	This enters play deactivated.
//	{Tap Effect}Destroy this. If you do, steal a soul from another player.
var momsShovel = engine.CardDef{
	Ref:               "moms_shovel",
	Kind:              engine.TreasureCard,
	Copies:            1,
	Soul:              1,
	Tap:               true,
	EntersDeactivated: true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Destroy this. If you do, steal a soul from another player.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.DestroyThis(engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					c.G.GainControl(c.Controller, othersSouls(c)[a[0]])
				}
			}, engine.Question{Text: "Steal which soul?", Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, id := range othersSouls(c) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			}}))},
		},
	},
}

func othersSouls(c *engine.Ctx) []engine.ObjectID {
	var out []engine.ObjectID
	for _, p := range otherPlayers(c) {
		out = append(out, souls(c.G, p)...)
	}
	return out
}
