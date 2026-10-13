package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pink Eye (Trinket Card)
//
//	Each time you take combat damage from a monster or player on an attack roll of 1, deal 1 damage to them.
//	-Trinket- This loot becomes an item under your control when it resolves.
var pinkEye = engine.CardDef{
	Ref:     "pink_eye",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you take combat damage from a monster or player on an attack roll of 1, deal 1 damage to them.",
		Trigger: combatDamageOnRollOf(1),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.DealDamageTo(engine.Target{Object: c.EventSource}, 1, c.Controller, c.Source)
		})},
	}},
}
