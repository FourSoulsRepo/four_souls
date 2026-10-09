package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Blind (Curse Card)
//
//	{Curse Effect}Monsters have +1{DC} on your turn.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfTheBlind = engine.CardDef{
	Ref:    "curse_of_the_blind",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Statics: []engine.Static{{Stat: engine.StatMonsterDC, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, _ engine.ObjectID) bool {
		return g.Turn.Active == g.Object(self).Controller
	}}},
}
