package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lazarus' Rags (Eternal Treasure Card)
//
//	Each time you die, after paying penalties, gain +1 treasure.
//	-Eternal- This can't be destroyed or put into discard.
var lazarusRags = engine.CardDef{
	Ref:     "lazarus_rags",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you die, after paying penalties, gain +1 treasure.",
			Trigger: engine.WhenYouDie(), // penalties are paid before triggers go on the stack
			Effects: []engine.Effect{engine.GainTreasure(1)},
		},
	},
}
