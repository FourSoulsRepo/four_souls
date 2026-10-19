package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Marbles (Active Treasure Card)
//
//	{Tap Effect}Choose a non-eternal card in play. Put a counter on it of a type already on it or remove a counter from it.
var marbles = engine.CardDef{
	Ref:    "marbles",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Choose a non-eternal card in play. Put a counter on it of a type already on it or remove a counter from it.",
		Costs: []engine.Cost{engine.Tap()},
		// Cards with counters are offered as it resolves: no target kind
		// covers every card in play.
		Effects: []engine.Effect{changeACounter(false)},
	}},
}
