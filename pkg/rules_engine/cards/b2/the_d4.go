package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The D4 (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose a player. They reroll each item they control.
var theD4 = engine.CardDef{
	Ref:    "the_d4",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Destroy this. If you do, choose a player. They reroll each item they control.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.DestroyThis(engine.EffectFunc(func(c *engine.Ctx) {
				for _, id := range itemsOf(c.G, c.Targets[0].Player) {
					c.G.RerollItem(id)
				}
			}))},
		},
	},
}
