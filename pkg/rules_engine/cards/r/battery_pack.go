package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Battery Pack (Paid Treasure Card)
//
//	At the end of your turn, put a counter on this.
//	{Paid Effect}Remove 2 counters from this:
//	Recharge each item a player controls.
var batteryPack = engine.CardDef{
	Ref:    "battery_pack",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, put a counter on this.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Remove 2 counters from this: Recharge each item a player controls.",
			Costs:   []engine.Cost{engine.RemoveCounters(2)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.RechargeItemsOf(0)},
		},
	},
}
