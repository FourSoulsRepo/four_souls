package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Blastocyst (Boss Card)
//
//	When this dies, expand monster slots by 2. The active player may attack the monster deck an additional time.
var blastocyst = engine.CardDef{
	Ref:     "blastocyst",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, expand monster slots by 2. The active player may attack the monster deck an additional time.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandMonsters(2); c.G.AddDeckAttack() })},
	}},
}
