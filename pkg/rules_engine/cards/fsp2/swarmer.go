package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Swarmer (Basic Monster Card)
//
//	When this dies, expand monster slots by 2.
var swarmer = engine.CardDef{
	Ref:     "swarmer",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, expand monster slots by 2.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandMonsters(2) })},
	}},
}
