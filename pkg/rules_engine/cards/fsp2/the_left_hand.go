package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Left Hand (Trinket Card)
//
//	Each time another player dies, you may recharge an item.
//	-Trinket- This loot becomes an item under your control when it resolves.
var theLeftHand = engine.CardDef{
	Ref:     "the_left_hand",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time another player dies, you may recharge an item.",
		Trigger: whenAnotherPlayerDies(),
		Effects: []engine.Effect{rechargeAnItem},
	}},
}
