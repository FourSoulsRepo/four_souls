package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Battery (Battery Card)
//
//	Choose a player. Recharge each item they control.
var megaBattery = engine.CardDef{
	Ref:    "mega_battery",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. Recharge each item they control.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.RechargeItemsOf(0)},
		},
	},
}
