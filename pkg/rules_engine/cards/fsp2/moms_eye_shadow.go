package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Eye Shadow (Passive Treasure Card)
//
//	If another player declares an attack on a monster, you may choose which monster they attack.
var momsEyeShadow = engine.CardDef{
	Ref:       "moms_eye_shadow",
	Kind:      engine.TreasureCard,
	Copies:    1,
	ThreePlus: true,
	Statics:   []engine.Static{engine.YouHave(engine.StatChooseAttack, 1)}, // "you may choose which monster they attack"
}
