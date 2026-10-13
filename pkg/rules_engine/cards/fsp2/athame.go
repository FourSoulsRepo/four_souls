package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Athame (Paid Treasure Card)
//
//	{Paid Effect}Discard a loot card:
//	Force a player to discard 2 loot cards.
var athame = engine.CardDef{
	Ref:    "athame",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "Discard a loot card: Force a player to discard 2 loot cards.",
		Costs:   []engine.Cost{engine.DiscardChosen(0)},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
		Effects: []engine.Effect{engine.Ask(forcedDiscard(2), forcedDiscardQuestions(2)...)},
	}},
}
