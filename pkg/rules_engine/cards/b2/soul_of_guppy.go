package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Guppy (Bonus Soul Card)
//
//	The first player to control 2 or more guppy items gains this soul.
var soulOfGuppy = engine.CardDef{
	Ref:       "soul_of_guppy",
	Kind:      engine.BonusSoulCard,
	Copies:    1,
	Soul:      1,
	BonusSoul: func(g *engine.Game, p engine.PlayerID) bool { return guppyItems(g, p) >= 2 },
}
