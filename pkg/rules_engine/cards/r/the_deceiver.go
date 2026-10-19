package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Deceiver (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theDeceiver = engine.CardDef{
	Ref:          "the_deceiver",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "ceremonial_blade",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
