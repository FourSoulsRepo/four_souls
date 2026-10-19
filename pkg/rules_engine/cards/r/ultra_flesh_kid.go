package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ultra Flesh Kid! (Soul Treasure Card)
//
//	Each time a monster or player dies, put a counter on this.
//	At the end of your turn, if this has 12+ counters, it becomes a soul and loses all abilities.
var ultraFleshKid = engine.CardDef{
	Ref:    "ultra_flesh_kid",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time a monster or player dies, put a counter on this.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool {
				return true
			}},
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind: engine.Triggered,
			Text: "At the end of your turn, if this has 12+ counters, it becomes a soul and loses all abilities.",
			Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				o := g.Object(self)
				return e.Player == o.Controller && o.CountersOf("") >= 12
			}},
			Effects: []engine.Effect{engine.BecomeSoul()},
		},
	},
}
