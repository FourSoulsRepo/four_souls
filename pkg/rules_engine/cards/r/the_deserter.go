package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Deserter (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theDeserter = engine.CardDef{
	Ref:          "the_deserter",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "anima_sola",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
