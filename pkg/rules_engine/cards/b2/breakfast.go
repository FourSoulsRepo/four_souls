package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Breakfast (Passive Treasure Card)
//
//	{HP}
var breakfast = engine.CardDef{
	Ref:     "breakfast",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerHP, 1)},
}
