package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Tiny Hands (Curse Card)
//
//	{Curse Effect}Your max hand size is 2.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfTinyHands = engine.CardDef{
	Ref:     "curse_of_tiny_hands",
	Kind:    engine.EventCard,
	Copies:  1,
	Curse:   true,
	Statics: []engine.Static{engine.YouHave(engine.StatHandCap, 2)},
}
