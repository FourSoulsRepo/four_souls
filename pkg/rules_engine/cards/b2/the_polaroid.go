package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Polaroid (Passive Treasure Card)
//
//	At the end of your turn, if you have 0 loot cards in your hand, loot 2.
var thePolaroid = engine.CardDef{
	Ref:    "the_polaroid",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, if you have 0 loot cards in your hand, loot 2.",
			Trigger: endOfTurnIf(func(g *engine.Game, p engine.PlayerID) bool { return len(g.Players[p].Hand) == 0 }),
			Effects: []engine.Effect{engine.Loot(2)},
		},
	},
}
