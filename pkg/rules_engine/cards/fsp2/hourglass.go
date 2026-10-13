package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hourglass (Passive Treasure Card)
//
//	Each time a player rolls a ❷, you may deactivate an item.
var hourglass = engine.CardDef{
	Ref:    "hourglass",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player rolls a 2, you may deactivate an item.",
		Trigger: engine.OnRollOf(2),
		Effects: []engine.Effect{deactivateAnItem},
	}},
}
