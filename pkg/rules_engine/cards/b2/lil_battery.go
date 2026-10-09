package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lil Battery (Battery Card)
//
//	Recharge an item.
var lilBattery = engine.CardDef{
	Ref:    "lil_battery",
	Kind:   engine.LootCard,
	Copies: 4,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Recharge an item.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetItem)},
			Effects: []engine.Effect{engine.RechargeTarget(0)},
		},
	},
}
