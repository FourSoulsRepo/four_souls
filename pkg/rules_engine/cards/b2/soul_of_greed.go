package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Greed (Bonus Soul Card)
//
//	The first player to have 25¢ or more gains this soul.
var soulOfGreed = engine.CardDef{
	Ref:       "soul_of_greed",
	Kind:      engine.BonusSoulCard,
	Copies:    1,
	Soul:      1,
	BonusSoul: func(g *engine.Game, p engine.PlayerID) bool { return g.Players[p].Cents >= 25 },
}
