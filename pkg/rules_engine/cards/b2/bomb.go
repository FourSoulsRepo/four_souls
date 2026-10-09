package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bomb! (Bomb Card)
//
//	Deal 1 damage to a Monster or Player.
var bomb = engine.CardDef{
	Ref:    "bomb",
	Kind:   engine.LootCard,
	Copies: 4,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Deal 1 damage to a monster or player.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
	},
}
