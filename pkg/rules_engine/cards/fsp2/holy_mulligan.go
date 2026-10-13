package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Mulligan (Holy/Charmed Monster Card)
//
//	When this dies, expand monster slots by 2. The active player may attack an additional time this turn.
var holyMulligan = engine.CardDef{
	Ref:     "holy_mulligan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, expand monster slots by 2. The active player may attack an additional time this turn.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandMonsters(2) }), engine.ActivePlayerAttacksAgain()},
	}},
}
