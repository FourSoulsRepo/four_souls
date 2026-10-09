package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sack Of Pennies (Active Treasure Card)
//
//	{Tap Effect}Gain 1¢.
//	Each time a player rolls a ❶, you may recharge this.
var sackOfPennies = engine.CardDef{
	Ref:    "sack_of_pennies",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Gain 1¢.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.GainCents(1)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 1, you may recharge this.",
			Trigger: engine.OnRollOf(1),
			Effects: []engine.Effect{may("Recharge Sack of Pennies?", engine.RechargeSelf())},
		},
	},
}
