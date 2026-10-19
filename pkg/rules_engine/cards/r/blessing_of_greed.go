package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of Greed (Room Card)
//
//	At the start of each turn, the active player gains 3¢.
var blessingOfGreed = engine.CardDef{
	Ref:    "blessing_of_greed",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of each turn, the active player gains 3¢.",
		Trigger: atStartOfTheTurn(),
		Effects: []engine.Effect{engine.GainCents(3)},
	}},
}
