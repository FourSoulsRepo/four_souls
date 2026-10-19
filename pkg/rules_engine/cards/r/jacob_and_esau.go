package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Jacob & Esau (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var jacobAndEsau = engine.CardDef{
	Ref:          "jacob_and_esau",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "sibling_rivalry",
	ThreePlus:    true,
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
