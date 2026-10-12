package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Placenta (Passive Treasure Card)
//
//	{HP}
//	You take no combat damage on attack rolls of 4.
var placenta = engine.CardDef{
	Ref:       "placenta",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Statics:   []engine.Static{engine.YouHave(engine.StatPlayerHP, 1)},
	CombatMod: noCombatDamageTakenOn(4),
}
