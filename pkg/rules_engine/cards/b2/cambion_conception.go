package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cambion Conception (Passive Treasure Card)
//
//	Each time you take damage, put counters on this equal to the amount of damage taken. Then, if this has 6+ counters, remove 6 counters from this and gain +1 treasure.
var cambionConception = engine.CardDef{
	Ref:    "cambion_conception",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, put counters on this equal to the amount of damage taken. Then, if this has 6+ counters, remove 6 counters from this and gain +1 treasure.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				o := c.G.Object(c.Source)
				if o.Zone.Kind != engine.ZoneInPlay {
					return
				}
				o.Counters = []engine.Counter{{Count: o.CountersOf("") + c.EventAmount}}
				if o.CountersOf("") >= 6 {
					o.Counters = []engine.Counter{{Count: o.CountersOf("") - 6}}
					c.Do(engine.GainTreasure(1))
				}
			})},
		},
	},
}
