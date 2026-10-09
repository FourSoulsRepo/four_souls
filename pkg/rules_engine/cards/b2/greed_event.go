package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greed! (Bad Event Card)
//
//	Choose the player with the most ¢ or tied for the most. That player loses all their ¢.
var greedEvent = engine.CardDef{
	Ref:    "greed_event",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Choose the player with the most ¢ or tied for the most. That player loses all their ¢.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					p := richest(c.G)[a[0]]
					c.Do(engine.EffectFunc(func(c *engine.Ctx) { c.G.Players[p].Cents = 0 }))
				}
			}, engine.Question{Text: "Which player loses all their ¢?", Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, p := range richest(c.G) {
					out = append(out, playerLabel(c.G, p))
				}
				return out
			}})},
		},
	},
}

// richest lists the players with the most ¢.
func richest(g *engine.Game) []engine.PlayerID {
	most := 0
	for _, pl := range g.Players {
		most = max(most, pl.Cents)
	}
	var out []engine.PlayerID
	for _, pl := range g.Players {
		if pl.Cents == most {
			out = append(out, pl.ID)
		}
	}
	return out
}
