package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: You gain +1{ATK} till the end of turn.
//	3-4: You gain +1{HP} till the end of turn.
//	5-6: Take 1 damage.
var pills3 = engine.CardDef{
	Ref:    "pills_3",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Roll- 1-2: You gain +1 ATK till the end of turn. 3-4: You gain +1 HP till the end of turn. 5-6: Take 1 damage.",
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.GainATKThisTurn(1, engine.You)).
				Results(3, 4, engine.GainHPThisTurn(1, engine.You)).
				Results(5, 6, engine.DealDamage(1, engine.You)))},
		},
	},
}
