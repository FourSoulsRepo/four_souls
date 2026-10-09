package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blue Baby (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var blueBaby = engine.CardDef{
	Ref:          "blue_baby",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "forever_alone",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
