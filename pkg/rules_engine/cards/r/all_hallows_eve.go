package r

import engine "github.com/FourSoulsRepo/rules_engine"

// All Hallows’ Eve (Room Card)
//
//	Players can't gain souls.
var allHallowsEve = engine.CardDef{
	Ref:     "all_hallows_eve",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{playersHave(engine.StatNoSouls, 1)},
}
