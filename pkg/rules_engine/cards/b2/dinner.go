package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dinner (Passive Treasure Card)
//
//	{HP}
var dinner = engine.CardDef{
	Ref:     "dinner",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerHP, 1)},
}
