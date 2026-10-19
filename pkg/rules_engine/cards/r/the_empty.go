package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Empty (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theEmpty = engine.CardDef{
	Ref:          "the_empty",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "abyss",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
