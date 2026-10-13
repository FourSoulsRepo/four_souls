package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum-Bo (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var bumBoFsp2 = engine.CardDef{
	Ref:          "bum_bo_fsp2",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "bag_o_trash",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
