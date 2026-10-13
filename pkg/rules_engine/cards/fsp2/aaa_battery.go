package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// AAA Battery (Trinket Card)
//
//	At the end of your turn, roll-
//	4-6: Recharge an item.
//	-Trinket- This loot becomes an item under your control when it resolves.
var aaaBattery = engine.CardDef{
	Ref:     "aaa_battery",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of your turn, roll- 4-6: Recharge an item.",
		Trigger: engine.AtEndOfYourTurn(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(4, 6, rechargeAnItem))},
	}},
}
