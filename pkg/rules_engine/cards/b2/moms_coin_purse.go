package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Coin Purse (Passive Treasure Card)
//
//	Loot +1 during your loot step.
var momsCoinPurse = engine.CardDef{
	Ref:     "moms_coin_purse",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatLootStep, 1)},
}
