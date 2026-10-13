package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Apollyon (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
var apollyon = engine.CardDef{
	Ref:          "apollyon",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "void",
	Tap:          true,
	Abilities:    []engine.Ability{extraLoot},
}
