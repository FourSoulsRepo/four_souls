package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Shop Upgrade! (Good Event Card)
//
//	Expand shop slots by 2.
//	The active player may attack an additional time this turn.
var shopUpgrade = engine.CardDef{
	Ref:    "shop_upgrade",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Expand shop slots by 2. The active player may attack an additional time this turn.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandShop(2) }), engine.AddAttacks(1, engine.You)},
		},
	},
}
