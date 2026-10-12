package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Azazel (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var azazel = engine.CardDef{
	Ref:          "azazel",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "lord_of_the_pit",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
