package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Enigma/amginE ehT (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	amginE ehT: Do not start the game with this side face up.
var theEnigma = engine.CardDef{
	Ref:          "the_enigma",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "flip",
	Tap:          true,
	Back:         "amgine_eht",
	Abilities:    []engine.Ability{extraLoot},
}

// amgineEht is The Enigma's back face (amginE ehT): it never starts
// face up.
var amgineEht = engine.CardDef{
	Ref:       "amgine_eht",
	Kind:      engine.CharacterCard,
	HP:        1,
	ATK:       2,
	Outside:   true,
	Tap:       true,
	Abilities: []engine.Ability{extraLoot},
}
