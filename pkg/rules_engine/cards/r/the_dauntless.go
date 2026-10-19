package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Dauntless (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theDauntless = engine.CardDef{
	Ref:          "the_dauntless",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           1,
	ATK:          1,
	StartingItem: "hypercoagulation",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
