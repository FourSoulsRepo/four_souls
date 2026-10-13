package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dark Judas (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var darkJudas = engine.CardDef{
	Ref:          "dark_judas",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "dark_arts",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
