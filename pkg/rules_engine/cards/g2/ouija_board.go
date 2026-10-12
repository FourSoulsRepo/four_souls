package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ouija Board (Active Treasure Card)
//
//	{Tap Effect}Reveal the top 4 cards of the monster deck. You may put an event card revealed this way in a monster slot not being attacked. Put the rest on the bottom of the deck in a random order.
var ouijaBoard = engine.CardDef{
	Ref:    "ouija_board",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Reveal the top 4 cards of the monster deck. You may put an event card revealed this way in a monster slot not being attacked. Put the rest on the bottom of the deck in a random order.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{ouija},
	}},
}
