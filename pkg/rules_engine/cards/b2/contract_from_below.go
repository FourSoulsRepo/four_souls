package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Contract From Below (Paid Treasure Card)
//
//	{Paid Effect}Destroy 2 items you control:
//	steal a non-eternal item from a player.
var contractFromBelow = engine.CardDef{
	Ref:    "contract_from_below",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Destroy 2 items you control: Steal a non-eternal item from a player.",
			Costs:   []engine.Cost{engine.DestroyChosen(0, 1)},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetYourItem, engine.NotChosen), engine.ChooseWhere(engine.TargetYourItem, engine.NotChosen), engine.ChooseWhere(engine.TargetItem, othersNonEternal)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.GainControl(c.Controller, c.Targets[2].Object) })},
		},
	},
}
