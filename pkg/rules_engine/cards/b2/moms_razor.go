package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Razor (Passive Treasure Card)
//
//	Each time a player rolls a ❻, you may deal 1 damage to them.
var momsRazor = engine.CardDef{
	Ref:    "moms_razor",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 6, you may deal 1 damage to them.",
			Trigger: engine.OnRollOf(6),
			Effects: []engine.Effect{may("Deal 1 damage to them?", engine.EffectFunc(func(c *engine.Ctx) {
				if c.EventPlayer != engine.NoPlayer {
					c.G.DealDamageTo(engine.Target{Player: c.EventPlayer, IsPlayer: true}, 1, c.Controller, c.Source)
				}
			}))},
		},
	},
}
