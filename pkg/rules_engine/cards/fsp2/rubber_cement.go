package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rubber Cement (Passive Treasure Card)
//
//	Each time you miss an attack roll, roll-
//	1-3: Deal 1 damage to a monster or player.
var rubberCement = engine.CardDef{
	Ref:    "rubber_cement",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you miss an attack roll, roll- 1-3: Deal 1 damage to a monster or player.",
		Trigger: whenYouMiss(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 3, damageAnything(1)))},
	}},
}
