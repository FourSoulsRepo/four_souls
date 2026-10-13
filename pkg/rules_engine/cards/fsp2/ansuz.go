package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ansuz (Pill/Rune Card)
//
//	Look at the top 4 cards of a deck and put them back in any order.
var ansuz = engine.CardDef{
	Ref:    "ansuz",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Look at the top 4 cards of a deck and put them back in any order.",
		Effects: []engine.Effect{reorderADeck(4)},
	}},
}
