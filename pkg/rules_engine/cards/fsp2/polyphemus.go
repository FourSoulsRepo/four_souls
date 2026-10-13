package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Polyphemus (Passive Treasure Card)
//
//	+2{ATK}
//	Monsters have +1{DC} on your turn.
var polyphemus = engine.CardDef{
	Ref:     "polyphemus",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerATK, 2), monstersDCOnYourTurn},
}
