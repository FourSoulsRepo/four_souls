package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Rib Of Greed (Trinket Card)
//
//	Abilities and the death penalty can't make you discard loot cards or lose ¢.
//	-Trinket- This loot becomes an item under your control when it resolves.
//
// Steals still work: a steal is giving, not losing (Q-17). Costs are paid
// by choice and still work too.
var ribOfGreed = engine.CardDef{
	Ref:     "rib_of_greed",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Statics: []engine.Static{engine.YouHave(engine.StatNoForcedLoss, 1)},
}
