package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Samson’s Blessing (Room Card)
//
//	Players have +1{ATK}.
var samsonsBlessing = engine.CardDef{
	Ref:     "samsons_blessing",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{playersHave(engine.StatPlayerATK, 1)},
}
