package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Book Of The Dead (Active Treasure Card)
//
//	{Tap Effect}Put a non-event monster card in discard in a monster slot not being attacked.
var bookOfTheDead = engine.CardDef{
	Ref:    "book_of_the_dead",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Put a non-event monster card in discard in a monster slot not being attacked.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{discardMonsterToSlot(1, false)},
	}},
}
