package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Heavy Is The Head (Room Card)
//
//	Players who control the most items or tied for the most may only recharge one item during their recharge step.
var heavyIsTheHead = engine.CardDef{
	Ref:     "heavy_is_the_head",
	Kind:    engine.RoomCard,
	Copies:  1,
	Statics: []engine.Static{playersWith(engine.StatRechargeLimit, 1, mostOf(itemCount))},
}
