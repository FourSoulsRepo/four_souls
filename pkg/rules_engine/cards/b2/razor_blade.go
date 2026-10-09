package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Razor Blade (Active Treasure Card)
//
//	{Tap Effect}Deal 1 damage to a player.
var razorBlade = engine.CardDef{
	Ref:    "razor_blade",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Deal 1 damage to a player.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
	},
}
