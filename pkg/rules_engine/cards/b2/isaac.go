package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Isaac (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var isaac = engine.CardDef{
	Ref:          "isaac",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "the_d6",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
