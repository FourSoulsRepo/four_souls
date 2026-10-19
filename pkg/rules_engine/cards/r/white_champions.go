package r

import engine "github.com/FourSoulsRepo/rules_engine"

// White Champions (Room Card)
//
//	Monsters have +1{DC}.
var whiteChampions = engine.CardDef{
	Ref:     "white_champions",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{engine.MonstersHave(engine.StatMonsterDC, 1)},
}
