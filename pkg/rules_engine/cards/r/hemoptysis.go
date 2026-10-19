package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hemoptysis (Eternal Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove a counter from this:
//	Choose a monster. Its {ATK} becomes 1.
//	{Paid Effect}Remove 2 counters from this:
//	Deal 1 damage to up to 2 monsters or players.
//	-Eternal- This can't be destroyed or put into discard.
var hemoptysis = engine.CardDef{
	Ref:     "hemoptysis",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Put a counter on this.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Remove a counter from this: Choose a monster. Its ATK becomes 1.",
			Costs:   []engine.Cost{engine.RemoveCounters(1)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonster)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Object(c.Targets[0].Object).ATKSet = 1 + 1
			})},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Remove 2 counters from this: Deal 1 damage to up to 2 monsters or players.",
			Costs:   []engine.Cost{engine.RemoveCounters(2)},
			Effects: []engine.Effect{damageUpTo(2, 1)},
		},
	},
}
