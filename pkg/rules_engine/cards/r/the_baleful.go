package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Baleful (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	When you control 1 or 2 souls, you have +1{ATK}. If you control 3+ souls, you have +2{ATK} instead.
var theBaleful = engine.CardDef{
	Ref:          "the_baleful",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           1,
	ATK:          1,
	StartingItem: "soulbond",
	ThreePlus:    true,
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
	Statics: []engine.Static{{
		Stat: engine.StatPlayerATK,
		Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
			return g.Object(self).Controller == p
		},
		Value: func(g *engine.Game, self engine.ObjectID) int {
			switch souls := g.SoulValue(g.Object(self).Controller); {
			case souls >= 3:
				return 2
			case souls >= 1:
				return 1
			}
			return 0
		},
	}},
}
