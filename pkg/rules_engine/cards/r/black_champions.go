package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Black Champions (Room Card)
//
//	Monsters have +1{ATK}.
var blackChampions = engine.CardDef{
	Ref:     "black_champions",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{engine.MonstersHave(engine.StatMonsterATK, 1)},
}
