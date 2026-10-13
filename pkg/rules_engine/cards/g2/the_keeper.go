package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Keeper (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theKeeper = engine.CardDef{
	Ref:          "the_keeper",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "wooden_nickel",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
