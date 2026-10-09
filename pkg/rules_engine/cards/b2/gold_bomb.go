package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gold Bomb!! (Bomb Card)
//
//	Deal 3 damage to a monster or player.
var goldBomb = engine.CardDef{
	Ref:    "gold_bomb",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Deal 3 damage to a monster or player.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.DealDamage(3, 0)},
		},
	},
}
