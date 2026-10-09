package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Tech X (Active Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove 3 counters from this:
//	Kill a player or monster.
var techX = engine.CardDef{
	Ref:    "tech_x",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Put a counter on this.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "Remove 3 counters from this: Kill a player or monster.",
			Costs:   []engine.Cost{engine.RemoveCounters(3)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.Kill(0)},
		},
	},
}
