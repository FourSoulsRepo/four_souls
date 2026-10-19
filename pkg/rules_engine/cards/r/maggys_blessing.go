package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Maggy’s Blessing (Room Card)
//
//	Players have +1{HP}.
var maggysBlessing = engine.CardDef{
	Ref:     "maggys_blessing",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{playersHave(engine.StatPlayerHP, 1)},
}
