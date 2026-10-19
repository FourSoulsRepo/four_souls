package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Hollow (Curse Card)
//
//	At the end of your turn, discard x loot cards and lose x cents, where x is the number of souls you control.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfTheHollow = engine.CardDef{
	Ref:    "curse_of_the_hollow",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of your turn, discard x loot cards and lose x cents, where x is the number of souls you control.",
		Trigger: engine.AtEndOfYourTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if x := c.G.SoulValue(c.Controller); x > 0 {
				c.Do(discardN(x), engine.LoseCents(x))
			}
		})},
	}},
}
