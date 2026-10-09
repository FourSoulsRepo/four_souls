package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Eye Of Greed (Passive Treasure Card)
//
//	Each time a player rolls a ❺, gain 3¢.
var eyeOfGreed = engine.CardDef{
	Ref:    "eye_of_greed",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 5, gain 3¢.",
			Trigger: engine.OnRollOf(5),
			Effects: []engine.Effect{engine.GainCents(3)},
		},
	},
}
