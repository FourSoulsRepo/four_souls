package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sumptorium (Eternal Treasure Card)
//
//	Each time you take damage, put a counter on this.
//	{Paid Effect}Remove 2 counters from this:
//	A monster gains -1{DC} till end of turn.
//	-Eternal- This can't be destroyed or put into discard.
var sumptorium = engine.CardDef{
	Ref:     "sumptorium",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, put a counter on this.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Remove 2 counters from this: A monster gains -1 DC till end of turn.",
			Costs:   []engine.Cost{engine.RemoveCounters(2)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonster)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatMonsterDC, Player: engine.NoPlayer, Object: c.Targets[0].Object, Amount: -1})
			})},
		},
	},
}
