package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Eve (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var eve = engine.CardDef{
	Ref:          "eve",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "the_curse",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
