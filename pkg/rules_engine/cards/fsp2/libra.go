package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Libra (Active Treasure Card)
//
//	{Tap Effect}Change the result of a dice roll to a 3.
var libra = engine.CardDef{
	Ref:    "libra",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Change the result of a dice roll to a 3.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if it, ok := c.G.StackItemByID(c.Targets[0].StackID); ok {
				c.G.SetRoll(it.ID, 3)
			}
		})},
	}},
}
