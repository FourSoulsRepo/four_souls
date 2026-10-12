package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Knife (Passive Treasure Card)
//
//	{ATK}
//	Each time you deal combat damage to a monster on an attack roll of 6, deal 1 damage to each other monster.
var momsKnife = engine.CardDef{
	Ref:     "moms_knife",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerATK, 1)},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you deal combat damage to a monster on an attack roll of 6, deal 1 damage to each other monster.",
		Trigger: combatDamageToAMonster(6),
		Effects: []engine.Effect{damageOtherMonsters},
	}},
}
