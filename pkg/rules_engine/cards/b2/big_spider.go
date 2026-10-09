package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Big Spider (Basic Monster Card)
//
//	When this dies, the active player may attack the monster deck an additional time.
var bigSpider = engine.CardDef{
	Ref:     "big_spider",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player may attack the monster deck an additional time.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.AddDeckAttack()
			})},
		},
	},
}
