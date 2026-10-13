package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Wiz (Passive Treasure Card)
//
//	Monsters have +1{DC} on your turn.
//	Each time you deal combat damage, deal 1 damage to another monster or player.
var theWiz = engine.CardDef{
	Ref:     "the_wiz",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{monstersDCOnYourTurn},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you deal combat damage, deal 1 damage to another monster or player.",
		Trigger: combatDamageToAMonster(),
		Effects: []engine.Effect{damageAnotherThan},
	}},
}
