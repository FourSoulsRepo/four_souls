package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Head Of Krampus (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1-3: Deal 1 damage to each player.
//	4-6: Deal 1 damage to each monster.
var headOfKrampus = engine.CardDef{
	Ref:    "head_of_krampus",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Roll- 1-3: Deal 1 damage to each player. 4-6: Deal 1 damage to each monster.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 3, engine.EachPlayer(engine.DealDamage(1, engine.You))).
			Results(4, 6, engine.EachMonsterTakesDamage(1)))},
	}},
}
