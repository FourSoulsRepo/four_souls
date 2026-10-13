package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Tape Worm (Trinket Card)
//
//	Each time you miss an attack roll, deal 1 damage to another player.
//	-Trinket- This loot becomes an item under your control when it resolves.
var tapeWorm = engine.CardDef{
	Ref:     "tape_worm",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you miss an attack roll, deal 1 damage to another player.",
		Trigger: whenYouMiss(),
		Effects: []engine.Effect{damageAPlayer(1, true)},
	}},
}
