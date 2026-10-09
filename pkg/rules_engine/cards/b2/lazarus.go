package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lazarus (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var lazarus = engine.CardDef{
	Ref:          "lazarus",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "lazarus_rags",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
