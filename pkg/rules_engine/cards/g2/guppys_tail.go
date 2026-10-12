package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Tail (Passive Treasure Card)
//
//	Loot +1 during your loot step.
//	At the end of your turn, discard a loot card.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
var guppysTail = engine.CardDef{
	Ref:     "guppys_tail",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Guppy:   true,
	Statics: []engine.Static{engine.YouHave(engine.StatLootStep, 1)},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of your turn, discard a loot card.",
		Trigger: engine.AtEndOfYourTurn(),
		Effects: []engine.Effect{discardOne},
	}},
}
