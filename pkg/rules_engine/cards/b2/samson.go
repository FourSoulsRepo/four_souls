package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Samson (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var samson = engine.CardDef{
	Ref:          "samson",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "blood_lust",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
