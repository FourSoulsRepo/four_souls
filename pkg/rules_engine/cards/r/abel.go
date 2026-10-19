package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Abel (Passive Treasure Card)
//
//	Each time you deal combat damage to a monster, roll-
//	4-6: Deal 1 damage to another monster or player.
var abel = engine.CardDef{
	Ref:       "abel",
	Kind:      engine.TreasureCard,
	Copies:    1,
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you deal combat damage to a monster, roll- 4-6: Deal 1 damage to another monster or player.",
		Trigger: combatDamageToAMonster(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(4, 6, damageOtherThanEvent))},
	}},
}
