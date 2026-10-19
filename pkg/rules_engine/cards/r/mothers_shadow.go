package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mother’s Shadow! (Bad Event Card)
//
//	The player who controls the most items or tied for the most dies. Put this in the monster deck 6 cards from the top.
var mothersShadow = engine.CardDef{
	Ref:    "mothers_shadow",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "The player who controls the most items or tied for the most dies. Put this in the monster deck 6 cards from the top.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			most := 0
			for _, pl := range c.G.Players {
				most = max(most, len(itemsOf(c.G, pl.ID)))
			}
			// Deaths go on the stack in reverse turn order, so they resolve
			// in turn order (R-MECH-28).
			order := c.G.TurnOrderFrom(c.G.Turn.Active)
			for i := len(order) - 1; i >= 0; i-- {
				if p := order[i]; len(itemsOf(c.G, p)) == most {
					them := *c
					them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: p}}
					them.Do(engine.Kill(0))
				}
			}
		}), backIntoDeck(6)},
	}},
}
