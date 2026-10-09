package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Trinity Shield (Passive Treasure Card)
//
//	Other players can't play loot cards or activate items on your turn.
var trinityShield = engine.CardDef{
	Ref:     "trinity_shield",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatLockOthers, 1)},
}
