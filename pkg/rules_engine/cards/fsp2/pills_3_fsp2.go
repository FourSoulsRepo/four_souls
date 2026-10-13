package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Deal 1 damage to a player.
//	3-4: Deal 1 damage to a monster.
//	5-6: Take 1 damage.
var pills3Fsp2 = engine.CardDef{
	Ref:    "pills_3_fsp2",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Roll- 1-2: Deal 1 damage to a player. 3-4: Deal 1 damage to a monster. 5-6: Take 1 damage.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, damageAPlayer(1, false)).
			Results(3, 4, damageAMonster(1)).
			Results(5, 6, engine.DealDamage(1, engine.You)))},
	}},
}
