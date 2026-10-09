package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Purse (Passive Treasure Card)
//
//	Loot +1 during your loot step.
var momsPurse = engine.CardDef{
	Ref:     "moms_purse",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatLootStep, 1)},
}
