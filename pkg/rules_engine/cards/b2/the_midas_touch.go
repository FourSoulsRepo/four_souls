package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Midas Touch (Passive Treasure Card)
//
//	Each time a monster dies, gain 3¢.
var theMidasTouch = engine.CardDef{
	Ref:    "the_midas_touch",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a monster dies, gain 3¢.",
			Trigger: engine.WhenAMonsterDies(),
			Effects: []engine.Effect{engine.GainCents(3)},
		},
	},
}
