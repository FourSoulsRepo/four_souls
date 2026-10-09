package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Charged Baby (Passive Treasure Card)
//
//	Each time a player rolls a ❷, you may recharge an item.
var chargedBaby = engine.CardDef{
	Ref:    "charged_baby",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 2, you may recharge an item.",
			Trigger: engine.OnRollOf(2),
			Effects: []engine.Effect{rechargeAnItem},
		},
	},
}
