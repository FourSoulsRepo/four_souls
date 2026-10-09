package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cain (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	If you control this as the game starts, you go first.
var cain = engine.CardDef{
	Ref:          "cain",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "sleight_of_hand",
	Tap:          true,
	GoesFirst:    true, // "If you control this as the game starts, you go first."
	Abilities:    []engine.Ability{extraLoot},
}
