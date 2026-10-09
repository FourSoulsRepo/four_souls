package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Paw (Active Treasure Card)
//
//	{Tap Effect}Pay 1{HP}. If you do, choose a player. Prevent the next instance of up to 2 damage they would take this turn.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
var guppysPaw = engine.CardDef{
	Ref:    "guppys_paw",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Guppy:  true,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Pay 1 HP. If you do, choose a player. Prevent the next instance of up to 2 damage they would take this turn.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if c.G.PayHP(c.Controller, 1) {
					c.Do(engine.PreventDamage(2, 0))
				}
			})},
		},
	},
}
