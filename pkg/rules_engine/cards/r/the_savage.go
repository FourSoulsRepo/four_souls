package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Savage (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn
var theSavage = engine.CardDef{
	Ref:          "the_savage",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "berserk",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
