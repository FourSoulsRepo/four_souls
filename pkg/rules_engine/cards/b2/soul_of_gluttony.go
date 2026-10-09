package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Gluttony (Bonus Soul Card)
//
//	The first player to have 10 or more loot cards in their hand gains this soul.
var soulOfGluttony = engine.CardDef{
	Ref:       "soul_of_gluttony",
	Kind:      engine.BonusSoulCard,
	Copies:    1,
	Soul:      1,
	BonusSoul: func(g *engine.Game, p engine.PlayerID) bool { return len(g.Players[p].Hand) >= 10 },
}
