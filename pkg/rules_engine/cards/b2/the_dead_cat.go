package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Dead Cat (Passive Treasure Card)
//
//	This item starts with 9 counters on it.
//	If you would take damage while this has counters on it, remove that many counters and prevent that much damage.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
var theDeadCat = engine.CardDef{
	Ref:                "the_dead_cat",
	Kind:               engine.TreasureCard,
	Copies:             1,
	EntersWithCounters: 9,
	DamageMod: func(g *engine.Game, self engine.ObjectID, t engine.Target, n int) int {
		o := g.Object(self)
		if !t.IsPlayer || t.Player != o.Controller || o.CountersOf("") == 0 {
			return n
		}
		used := min(n, o.CountersOf(""))
		o.Counters = []engine.Counter{{Count: o.CountersOf("") - used}}
		return n - used
	},
}
