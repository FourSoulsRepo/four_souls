package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Maggy (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var maggy = engine.CardDef{
	Ref:          "maggy",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "yum_heart",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
