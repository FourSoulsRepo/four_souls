package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lard (Passive Treasure Card)
//
//	+2{HP}
//	Each time you take damage, discard a loot card.
var lard = engine.CardDef{
	Ref:     "lard",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerHP, 2)},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you take damage, discard a loot card.",
		Trigger: engine.WhenYouTakeDamage(),
		Effects: []engine.Effect{discardOne},
	}},
}
