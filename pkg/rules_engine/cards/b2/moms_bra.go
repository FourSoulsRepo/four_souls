package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Bra (Active Treasure Card)
//
//	{Tap Effect}Choose a monster or player. The next instance of damage they take this turn is reduced to 1.
var momsBra = engine.CardDef{
	Ref:    "moms_bra",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a monster or player. The next instance of damage they take this turn is reduced to 1.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.CapNextDamage(1, 0)},
		},
	},
}
