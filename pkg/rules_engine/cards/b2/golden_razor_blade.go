package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Razor Blade (Paid Treasure Card)
//
//	{Paid Effect}Pay 5¢:
//	Deal 1 damage to a monster or player.
var goldenRazorBlade = engine.CardDef{
	Ref:    "golden_razor_blade",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Pay 5¢: Deal 1 damage to a monster or player.",
			Costs:   []engine.Cost{engine.PayCents(5)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
	},
}
