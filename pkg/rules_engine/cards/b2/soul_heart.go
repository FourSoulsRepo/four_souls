package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Heart (Dice Shard/Soul Heart Card)
//
//	Choose a player. Prevent the next 1 damage they would take this turn.
var soulHeart = engine.CardDef{
	Ref:    "soul_heart",
	Kind:   engine.LootCard,
	Copies: 2,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. Prevent the next 1 damage they would take this turn.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.PreventDamage(1, 0)},
		},
	},
}
