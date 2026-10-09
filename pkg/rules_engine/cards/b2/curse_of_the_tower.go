package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Tower (Passive Treasure Card)
//
//	Each time you take damage, roll-
//	1-3: Each other player takes 1 damage.
//	4-6: Deal 1 damage to a monster.
var curseOfTheTower = engine.CardDef{
	Ref:    "curse_of_the_tower",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, roll- 1-3: Each other player takes 1 damage. 4-6: Deal 1 damage to a monster.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 3, engine.EachOtherPlayer(engine.DealDamage(1, engine.You))).Results(4, 6, damageAMonster(1)))},
		},
	},
}
