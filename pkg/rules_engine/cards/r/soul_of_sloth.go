package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Sloth (Bonus Soul Card)
//
//	The first time a player controls 4 items, the active player chooses a player who controls the fewest items or tied for fewest. That player gains this soul.
var soulOfSloth = engine.CardDef{
	Ref:    "soul_of_sloth",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "The first time a player controls 4 items, the active player chooses a player who controls the fewest items or tied for fewest. That player gains this soul.",
			Trigger: fourItems(engine.EvGainedTreasure, "stolen"),
			Effects: []engine.Effect{fewestGetsThis(func(g *engine.Game, p engine.PlayerID) int { return len(itemsOf(g, p)) })},
		},
		{
			Kind:    engine.Triggered,
			Text:    "The first time a player controls 4 items, the active player chooses a player who controls the fewest items or tied for fewest. That player gains this soul.",
			Trigger: fourItems(engine.EvGainedItem, ""),
			Effects: []engine.Effect{fewestGetsThis(func(g *engine.Game, p engine.PlayerID) int { return len(itemsOf(g, p)) })},
		},
	},
}

// fourItems triggers when a player gains or steals an item and now
// controls 4, if this hasn't happened yet.
func fourItems(on engine.EventKind, text string) engine.Trigger {
	return engine.Trigger{On: on, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == text && g.Object(self).CountersOf(firstTime) == 0 && e.Player != engine.NoPlayer && len(itemsOf(g, e.Player)) >= 4
	}}
}
