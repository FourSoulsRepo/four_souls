package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Whore Of Babylon (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var whoreOfBabylon = engine.CardDef{
	Ref:          "whore_of_babylon",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "gimpy",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
