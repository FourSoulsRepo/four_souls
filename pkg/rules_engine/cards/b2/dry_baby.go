package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dry Baby (Passive Treasure Card)
//
//	Damage you would take is reduced to 1.
var dryBaby = engine.CardDef{
	Ref:    "dry_baby",
	Kind:   engine.TreasureCard,
	Copies: 1,
	DamageMod: func(g *engine.Game, self engine.ObjectID, t engine.Target, n int) int {
		if t.IsPlayer && t.Player == g.Object(self).Controller {
			return min(n, 1)
		}
		return n
	},
}
