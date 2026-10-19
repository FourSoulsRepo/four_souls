package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Judas’ Shadow (Passive Treasure Card)
//
//	Each time you die, before paying penalties, put 2 counters on this.
//	When this is destroyed, gain x¢ and loot x, where x is equal to the number of counters on this.
var judasShadow = engine.CardDef{
	Ref:    "judas_shadow",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you die, before paying penalties, put 2 counters on this.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{engine.AddCounters(2)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this is destroyed, gain x¢ and loot x, where x is equal to the number of counters on this.",
			Trigger: engine.WhenThisIsDestroyed(),
			// The destroyed object keeps the counters it had.
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if x := c.G.Object(c.Source).CountersOf(""); x > 0 {
					c.Do(engine.GainCents(x), engine.Loot(x))
				}
			})},
		},
	},
}
