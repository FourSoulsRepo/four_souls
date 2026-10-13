package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Infestation (Eternal Treasure Card)
//
//	{Tap Effect}Loot 2, then discard 1 loot card.
//	-Eternal- This can't be destroyed or put into discard.
var infestation = engine.CardDef{
	Ref:     "infestation",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Loot 2, then discard 1 loot card.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.Loot(2), discardOne},
	}},
}
