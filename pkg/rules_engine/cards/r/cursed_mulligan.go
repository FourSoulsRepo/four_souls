package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Mulligan (Cursed Monster Card)
//
//	Each time a monster dies, expand monster slots by 1.
var cursedMulligan = engine.CardDef{
	Ref:     "cursed_mulligan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a monster dies, expand monster slots by 1.",
		Trigger: engine.WhenAMonsterDies(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.ExpandMonsters(1) })},
	}},
}
