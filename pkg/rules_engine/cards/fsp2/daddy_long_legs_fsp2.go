package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Daddy Long Legs (Passive Treasure Card)
//
//	At the start of your turn, deal 1 damage to a monster or player.
var daddyLongLegsFsp2 = engine.CardDef{
	Ref:    "daddy_long_legs_fsp2",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn, deal 1 damage to a monster or player.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{damageAnything(1)},
	}},
}
