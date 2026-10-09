package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pay To Play (Paid Treasure Card)
//
//	{Paid Effect}Pay 10¢:
//	Steal a non-eternal item a player controls.
var payToPlay = engine.CardDef{
	Ref:    "pay_to_play",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Pay 10¢: Steal a non-eternal item a player controls.",
			Costs:   []engine.Cost{engine.PayCents(10)},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, othersNonEternal)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.GainControl(c.Controller, c.Targets[0].Object)
			})},
		},
	},
}
