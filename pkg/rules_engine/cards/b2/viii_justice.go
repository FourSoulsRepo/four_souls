package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// VIII. Justice (Wildcard Card)
//
//	Choose a player. Loot and gain ¢ until you have the same number of each as they do.
var viiiJustice = engine.CardDef{
	Ref:    "viii_justice",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose a player. Loot and gain ¢ until you have the same number of each as they do.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				me, them := c.G.Players[c.Controller], c.G.Players[c.Targets[0].Player]
				if n := len(them.Hand) - len(me.Hand); n > 0 {
					c.Do(engine.Loot(n))
				}
				if n := them.Cents - me.Cents; n > 0 {
					c.Do(engine.GainCents(n))
				}
			})},
		},
	},
}
