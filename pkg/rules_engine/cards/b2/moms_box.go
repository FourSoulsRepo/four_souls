package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Box (Passive Treasure Card)
//
//	Each time a player rolls a ❹, you may loot 1, then discard a loot card.
var momsBox = engine.CardDef{
	Ref:    "moms_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 4, you may loot 1, then discard a loot card.",
			Trigger: engine.OnRollOf(4),
			Effects: []engine.Effect{may("Loot 1, then discard a loot card?", engine.Loot(1), discardOne)},
		},
	},
}
