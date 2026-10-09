package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Compost (Active Treasure Card)
//
//	{Tap Effect}The next time a player would loot, they loot from the top of the loot discard instead.
var compost = engine.CardDef{
	Ref:    "compost",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: The next time a player would loot, they loot from the top of the loot discard instead.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.Compost = true })},
		},
	},
}
