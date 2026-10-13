package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Recharge each item you control.
//	3-4: Monsters you attack have -1{DC} till end of turn.
//	5-6: Monsters you attack have +2{DC} till end of turn.
var pillsG2 = engine.CardDef{
	Ref:    "pills_g2",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Roll- 1-2: Recharge each item you control. 3-4: Monsters you attack have -1 DC till end of turn. 5-6: Monsters you attack have +2 DC till end of turn.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, rechargeYourItems).
			Results(3, 4, attackedDC(-1)).
			Results(5, 6, attackedDC(2)))},
	}},
}
