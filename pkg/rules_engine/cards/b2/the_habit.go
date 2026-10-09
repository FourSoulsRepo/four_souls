package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Habit (Passive Treasure Card)
//
//	The first time you take damage each turn, you may recharge an item.
var theHabit = engine.CardDef{
	Ref:    "the_habit",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "The first time you take damage each turn, you may recharge an item.",
			Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				p := g.Object(self).Controller
				return e.Player == p && e.Object == 0 && g.Players[p].TimesDamaged == 1
			}},
			Effects: []engine.Effect{rechargeAnItem},
		},
	},
}
