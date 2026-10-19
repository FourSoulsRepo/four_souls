package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Handicapped Placard (Passive Treasure Card)
//
//	Monster have -x{DC} on your turn, where x is the number of souls the player with the most souls controls minus the number of souls you control.
var handicappedPlacard = engine.CardDef{
	Ref:    "handicapped_placard",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Statics: []engine.Static{{
		Stat: engine.StatMonsterDC,
		Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, _ engine.ObjectID) bool {
			return g.Turn.Active == g.Object(self).Controller
		},
		Value: func(g *engine.Game, self engine.ObjectID) int { return -soulsBehind(g, g.Object(self).Controller) },
	}},
}
