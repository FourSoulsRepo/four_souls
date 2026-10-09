package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mr. Boom (Active Treasure Card)
//
//	{Tap Effect}Deal 1 damage to a monster.
var mrBoom = engine.CardDef{
	Ref:    "mr_boom",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Deal 1 damage to a monster.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonster)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
	},
}
