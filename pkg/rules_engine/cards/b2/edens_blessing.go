package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Eden’s Blessing (Passive Treasure Card)
//
//	At the end of your turn, if you have 0¢, gain 6¢.
var edensBlessing = engine.CardDef{
	Ref:    "edens_blessing",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, if you have 0¢, gain 6¢.",
			Trigger: endOfTurnIf(func(g *engine.Game, p engine.PlayerID) bool { return g.Players[p].Cents == 0 }),
			Effects: []engine.Effect{engine.GainCents(6)},
		},
	},
}
