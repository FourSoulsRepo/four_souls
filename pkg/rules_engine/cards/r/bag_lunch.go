package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bag Lunch (Trinket Card)
//
//	If this would be destroyed, if it has no counters on it, put a counter on it instead.
//	You have +1{HP} while this has a counter on it.
//	-Trinket- This loot becomes an item under your control when it resolves.
var bagLunch = engine.CardDef{
	Ref:     "bag_lunch",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	InsteadOfDestroy: func(g *engine.Game, self engine.ObjectID) bool {
		o := g.Object(self)
		if o.Role != engine.RoleItem || len(o.Counters) > 0 {
			return false
		}
		g.AddCountersTo(self, 1)
		return true
	},
	Statics: []engine.Static{{Stat: engine.StatPlayerHP, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		o := g.Object(self)
		return o.Controller == p && len(o.Counters) > 0
	}}},
}
