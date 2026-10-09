package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Starter Deck (Passive Treasure Card)
//
//	At the end of your turn, if you have 8 or more loot cards in your hand, loot 2.
var starterDeck = engine.CardDef{
	Ref:    "starter_deck",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, if you have 8 or more loot cards in your hand, loot 2.",
			Trigger: endOfTurnIf(func(g *engine.Game, p engine.PlayerID) bool { return len(g.Players[p].Hand) >= 8 }),
			Effects: []engine.Effect{engine.Loot(2)},
		},
	},
}
