package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Magnet (Passive Treasure Card)
//
//	Each time another player gains ¢, they must give you 1¢.
var magnet = engine.CardDef{
	Ref:    "magnet",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time another player gains ¢, they must give you 1¢.",
		Trigger: engine.Trigger{On: engine.EvGainedCents, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Player != engine.NoPlayer && e.Player != g.Object(self).Controller && e.Amount > 0
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.GiveCents(c.EventPlayer, c.Controller, 1)
		})},
	}},
}
