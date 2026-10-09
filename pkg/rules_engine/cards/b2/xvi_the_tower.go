package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XVI. The Tower (Wildcard Card)
//
//	Roll-
//	1-2: Each player takes 1 damage.
//	3-4: Each monster takes 1 damage.
//	5-6: Each player takes 2 damage.
var xviTheTower = engine.CardDef{
	Ref:    "xvi_the_tower",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Roll- 1-2: Each player takes 1 damage. 3-4: Each monster takes 1 damage. 5-6: Each player takes 2 damage.",
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, engine.EachPlayer(engine.DealDamage(1, engine.You))).
				Results(3, 4, engine.EachMonsterTakesDamage(1)).
				Results(5, 6, engine.EachPlayer(engine.DealDamage(2, engine.You))))},
		},
	},
}
