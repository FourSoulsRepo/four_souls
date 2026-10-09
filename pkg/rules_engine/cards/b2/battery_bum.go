package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Battery Bum (Paid Treasure Card)
//
//	{Paid Effect}Pay 4¢:
//	Recharge an item.
var batteryBum = engine.CardDef{
	Ref:    "battery_bum",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Pay 4¢: Recharge an item.",
			Costs:   []engine.Cost{engine.PayCents(4)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetItem)},
			Effects: []engine.Effect{engine.RechargeTarget(0)},
		},
	},
}
