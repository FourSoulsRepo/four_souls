package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XXI. The World (Wildcard Card)
//
//	Look at each player's hand, then loot 2.
var xxiTheWorld = engine.CardDef{
	Ref:    "xxi_the_world",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Look at each player's hand, then loot 2.",
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				for _, pl := range c.G.Players {
					if pl.ID != c.Controller {
						c.G.LookAt(c.Controller, pl.Hand...)
					}
				}
				c.Do(engine.Loot(2))
			})},
		},
	},
}
