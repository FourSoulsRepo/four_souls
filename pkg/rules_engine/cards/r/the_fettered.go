package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Fettered (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theFettered = engine.CardDef{
	Ref:          "the_fettered",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "dead_weight",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
