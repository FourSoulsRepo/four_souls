package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Abaddon (Passive Treasure Card)
//
//	+3{ATK}
//	Each time you take damage, die!
var abaddon = engine.CardDef{
	Ref:     "abaddon",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerATK, 3)},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you take damage, die!",
		Trigger: engine.WhenYouTakeDamage(),
		Effects: []engine.Effect{killYou},
	}},
}
