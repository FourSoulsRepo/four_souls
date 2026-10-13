package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Globin (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➂, each monster heals 2{HP}.
var cursedGlobin = engine.CardDef{
	Ref:     "cursed_globin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player rolls a 3, each monster heals 2 HP.",
		Trigger: engine.OnRollOf(3),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for _, id := range monstersInPlay(c.G) {
				c.G.HealObject(id, 2)
			}
		})},
	}},
}
