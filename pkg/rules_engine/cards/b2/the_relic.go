package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Relic (Passive Treasure Card)
//
//	Each time a player rolls a ❶, loot 1.
var theRelic = engine.CardDef{
	Ref:    "the_relic",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 1, loot 1.",
			Trigger: engine.OnRollOf(1),
			Effects: []engine.Effect{engine.Loot(1)},
		},
	},
}
