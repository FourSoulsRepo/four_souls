package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ball Of Bandages (Passive Treasure Card)
//
//	You take no combat damage on attack rolls of 3.
var ballOfBandages = engine.CardDef{
	Ref:       "ball_of_bandages",
	Kind:      engine.TreasureCard,
	Copies:    1,
	CombatMod: noCombatDamageTakenOn(3),
}
