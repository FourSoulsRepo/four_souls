package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Eternal D6 (Active Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove 2 counters from this:
//	Till end of turn, if a player would roll a 2 or 3, it becomes a 1 instead.
var eternalD6 = engine.CardDef{
	Ref:    "eternal_d6",
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
			Kind:  engine.Activated,
			Text:  "$ Remove 2 counters from this: Till end of turn, if a player would roll a 2 or 3, it becomes a 1 instead.",
			Costs: []engine.Cost{engine.RemoveCounters(2)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatRollToOne, Player: engine.NoPlayer, Amount: 1})
			})},
		},
	},
}
