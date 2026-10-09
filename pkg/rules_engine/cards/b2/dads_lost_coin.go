package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dad’s Lost Coin (Passive Treasure Card)
//
//	Each time a player would roll a ❶, you may force that player to reroll it.
var dadsLostCoin = engine.CardDef{
	Ref:    "dads_lost_coin",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player would roll a 1, you may force that player to reroll it.",
			Trigger: engine.WhenARollWouldBe(1),
			Effects: []engine.Effect{may("Force a reroll?", engine.EffectFunc(func(c *engine.Ctx) { c.G.Reroll(c.EventStack) }))},
		},
	},
}
