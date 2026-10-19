package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Champions (Room Card)
//
//	Monsters have +1{HP}.
var redChampions = engine.CardDef{
	Ref:     "red_champions",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{engine.MonstersHave(engine.StatMonsterHP, 1)},
}
