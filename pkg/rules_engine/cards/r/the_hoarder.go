package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Hoarder (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theHoarder = engine.CardDef{
	Ref:          "the_hoarder",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "bag_of_crafting",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
