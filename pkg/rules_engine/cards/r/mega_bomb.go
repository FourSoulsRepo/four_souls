package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Bomb! (Bomb Card)
//
//	Deal 2 damage to a monster or player.
var megaBomb = engine.CardDef{
	Ref:    "mega_bomb",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Deal 2 damage to a monster or player.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
		Effects: []engine.Effect{engine.DealDamage(2, 0)},
	}},
}
