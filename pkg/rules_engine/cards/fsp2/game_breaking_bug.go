package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Game Breaking Bug! (Passive Treasure Card)
//
//	Each time a player rolls a ❶, you may reroll an item they control.
var gameBreakingBug = engine.CardDef{
	Ref:    "game_breaking_bug",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player rolls a 1, you may reroll an item they control.",
		Trigger: engine.OnRollOf(1),
		Effects: []engine.Effect{rerollRollersItem},
	}},
}
