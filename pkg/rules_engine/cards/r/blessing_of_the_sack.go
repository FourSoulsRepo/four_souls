package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blessing Of The Sack (Room Card)
//
//	Rewards are doubled.
var blessingOfTheSack = engine.CardDef{
	Ref:     "blessing_of_the_sack",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{engine.MonstersHave(engine.StatDoubleRewards, 1)},
}
