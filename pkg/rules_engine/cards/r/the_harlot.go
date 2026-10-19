package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Harlot (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theHarlot = engine.CardDef{
	Ref:          "the_harlot",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "gello",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
