package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Judas (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var judas = engine.CardDef{
	Ref:          "judas",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "book_of_belial",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
