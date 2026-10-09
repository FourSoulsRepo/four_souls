package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Loss (Curse Card)
//
//	{Curse Effect}When you die, destroy a soul you control.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfLoss = engine.CardDef{
	Ref:    "curse_of_loss",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you die, destroy a soul you control.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					c.G.DestroyObject(c.Controller, souls(c.G, c.Controller)[a[0]])
				}
			}, engine.Question{Text: "Destroy which soul?", Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, id := range souls(c.G, c.Controller) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			}})},
		},
	},
}
