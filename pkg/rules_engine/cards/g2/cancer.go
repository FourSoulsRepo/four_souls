package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cancer (Trinket Card)
//
//	You have +1 to your first dice roll each turn.
//	-Trinket- This loot becomes an item under your control when it resolves.
var cancer = engine.CardDef{
	Ref:     "cancer",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Statics: []engine.Static{{
		Stat:   engine.StatRoll,
		Amount: 1,
		Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
			return g.Object(self).Controller == p && g.Players[p].Rolls == 0 // "your first dice roll each turn"
		},
	}},
}
