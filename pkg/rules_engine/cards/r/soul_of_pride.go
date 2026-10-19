package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Soul Of Pride (Bonus Soul Card)
//
//	Each time a player gains a treasure, put a counter on this.
//	When the 6th counter is put on this, the active player gains this soul.
var soulOfPride = engine.CardDef{
	Ref:    "soul_of_pride",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player gains a treasure, put a counter on this.",
			Trigger: cardkit.WhenAPlayerGainsATreasure(),
			Effects: []engine.Effect{counterOnBonusSoul},
		},
		sixthCounterGivesSoul,
	},
}
