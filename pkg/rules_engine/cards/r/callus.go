package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Callus (Trinket Card)
//
//	Prevent all non-combat damage you would take.
//	-Trinket- This loot becomes an item under your control when it resolves.
var callus = engine.CardDef{
	Ref:     "callus",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	DamageMod: func(g *engine.Game, self engine.ObjectID, t engine.Target, n int, combat bool) int {
		if o := g.Object(self); !combat && t.IsPlayer && t.Player == o.Controller && o.Role == engine.RoleItem {
			return 0
		}
		return n
	},
}
