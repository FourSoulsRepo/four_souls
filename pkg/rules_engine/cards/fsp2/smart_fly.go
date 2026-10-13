package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Smart Fly (Active Treasure Card)
//
//	{Tap Effect}Look at the top card of a deck. You may put it into discard or put it back on top.
//	Each time you take damage, you may recharge this.
var smartFly = engine.CardDef{
	Ref:    "smart_fly",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Look at the top card of a deck. You may put it into discard or put it back on top.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{lookMayDiscard},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, you may recharge this.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{may("Recharge Smart Fly?", engine.RechargeSelf())},
		},
	},
}
