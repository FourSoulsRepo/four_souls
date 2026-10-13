package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var guppy = engine.CardDef{
	Ref:          "guppy",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "infestation",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
