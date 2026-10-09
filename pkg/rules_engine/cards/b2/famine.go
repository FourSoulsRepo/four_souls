package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Famine (Boss Card)
//
//	When this dies, the active player skips their next turn.
var famine = engine.CardDef{
	Ref:     "famine",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player skips their next turn.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Players[c.Controller].SkipTurns++
			})},
		},
	},
}
