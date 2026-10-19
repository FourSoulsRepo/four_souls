package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sacrificial Dagger (Paid Treasure Card)
//
//	Each time a player rolls a ❶, put a counter on this.
//	{Paid Effect}Remove 3 counters from this:
//	Deal 1 damage to a monster or player.
var sacrificialDagger = engine.CardDef{
	Ref:    "sacrificial_dagger",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a ❶, put a counter on this.",
			Trigger: engine.OnRollOf(1),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Remove 3 counters from this: Deal 1 damage to a monster or player.",
			Costs:   []engine.Cost{engine.RemoveCounters(3)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
	},
}
