package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fruitcake (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1: Loot 1.
//	2-3: Gain 1¢.
//	4-6: Look at the top card of a deck.
//	Each time a player rolls the same result twice in a row on the same turn, recharge this.
var fruitcake = engine.CardDef{
	Ref:    "fruitcake",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Roll- 1: Loot 1. 2-3: Gain 1¢. 4-6: Look at the top card of a deck.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 1, engine.Loot(1)).
				Results(2, 3, engine.GainCents(1)).
				Results(4, 6, lookAtADeck))},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls the same result twice in a row on the same turn, recharge this.",
			Trigger: engine.Trigger{On: engine.EvRollRepeated, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }}, // R-DICE-08, Q-10
			Effects: []engine.Effect{engine.RechargeSelf()},
		},
	},
}
