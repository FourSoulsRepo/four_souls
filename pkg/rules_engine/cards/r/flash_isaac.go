package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Flash Isaac (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var flashIsaac = engine.CardDef{
	Ref:          "flash_isaac",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "classic_roller",
	ThreePlus:    true,
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
