package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Soiled (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theSoiled = engine.CardDef{
	Ref:          "the_soiled",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "ibs",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
