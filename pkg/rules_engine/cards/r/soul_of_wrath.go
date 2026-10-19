package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Wrath (Bonus Soul Card)
//
//	Each time a player dies, put a counter on this.
//	When the 6th counter is put on this, the active player gains this soul.
var soulOfWrath = engine.CardDef{
	Ref:    "soul_of_wrath",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player dies, put a counter on this.",
			Trigger: engine.WhenAPlayerDies(),
			Effects: []engine.Effect{counterOnBonusSoul},
		},
		sixthCounterGivesSoul,
	},
}
