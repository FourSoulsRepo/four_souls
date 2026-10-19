package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Broken (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var theBroken = engine.CardDef{
	Ref:          "the_broken",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "spindown_dice",
	ThreePlus:    true,
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
