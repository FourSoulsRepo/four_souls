package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Box! (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, you may play any number of additional loot cards till end of turn.
var box = engine.CardDef{
	Ref:    "box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Destroy this. If you do, you may play any number of additional loot cards till end of turn.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.DestroyThis(engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Players[c.Controller].ExtraLootPlays += 99 // any number
			}))},
		},
	},
}
