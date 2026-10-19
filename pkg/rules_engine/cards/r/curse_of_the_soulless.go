package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Soulless (Curse Card)
//
//	You can't gain souls.
//	If this would be put into discard, instead give it to another player.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfTheSoulless = engine.CardDef{
	Ref:     "curse_of_the_soulless",
	Kind:    engine.EventCard,
	Copies:  1,
	Curse:   true,
	Statics: []engine.Static{engine.YouHave(engine.StatNoSouls, 1)},
	// "If this would be put into discard, instead give it to another
	// player": when destroyed, or when its holder dies.
	InsteadOfDestroy: func(g *engine.Game, self engine.ObjectID) bool {
		if g.Object(self).Role != engine.RoleCurse {
			return false
		}
		g.PassCurse(self)
		return true
	},
}
