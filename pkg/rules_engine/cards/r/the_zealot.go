package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Zealot (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theZealot = engine.CardDef{
	Ref:          "the_zealot",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "lemegeton",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
