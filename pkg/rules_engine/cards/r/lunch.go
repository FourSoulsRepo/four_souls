package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lunch (Passive Treasure Card)
//
//	+1{HP}
var lunch = engine.CardDef{
	Ref:     "lunch",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerHP, 1)},
}
