package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Synthoil (Passive Treasure Card)
//
//	You have +1 to attack rolls.
var synthoil = engine.CardDef{
	Ref:     "synthoil",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatAttackRoll, 1)},
}
