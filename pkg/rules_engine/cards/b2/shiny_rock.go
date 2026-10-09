package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Shiny Rock (Passive Treasure Card)
//
//	Each time you activate an item, gain 1¢.
var shinyRock = engine.CardDef{
	Ref:    "shiny_rock",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time you activate an item, gain 1¢.",
			Trigger: engine.Trigger{On: engine.EvActivated, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Player == g.Object(self).Controller && g.Object(e.Object).Role == engine.RoleItem
			}},
			Effects: []engine.Effect{engine.GainCents(1)},
		},
	},
}
