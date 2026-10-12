package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Charged Penny (1 Cent Card)
//
//	Gain 1¢, then recharge an item.
var chargedPenny = engine.CardDef{
	Ref:       "charged_penny",
	Kind:      engine.LootCard,
	Copies:    1,
	Abilities: []engine.Ability{{Kind: engine.LootAbility, Text: "Gain 1¢, then recharge an item.", Effects: []engine.Effect{engine.GainCents(1), rechargeAnItem}}},
}
