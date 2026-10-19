package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll 3 times!-
//	1-2: You gain +1{ATK} till end of turn!
//	3-4: You gain +1{HP} till end of turn!
//	5-6: You take 1 damage!
var pills2R = engine.CardDef{
	Ref:    "pills_2_r",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Roll 3 times!- 1-2: You gain +1 ATK till end of turn! 3-4: You gain +1 HP till end of turn! 5-6: You take 1 damage!",
		Effects: []engine.Effect{engine.Roll(pills2Table), engine.Roll(pills2Table), engine.Roll(pills2Table)},
	}},
}

var pills2Table = engine.RollTable{}.
	Results(1, 2, engine.GainATKThisTurn(1, engine.You)).
	Results(3, 4, engine.GainHPThisTurn(1, engine.You)).
	Results(5, 6, engine.DealDamage(1, engine.You))
