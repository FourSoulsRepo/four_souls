package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Rag Mega (Boss Card)
//
//	When an attack is declared on this, each non-active player discards a loot card.
//	When this dies, the active player kills up to 3 other players.
var ragMega = engine.CardDef{
	Ref:     "rag_mega",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When an attack is declared on this, each non-active player discards a loot card.",
			Trigger: whenAttackedThis(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				for _, p := range othersOfActive(c.G) {
					them := *c
					them.Controller = p
					them.Do(discardOne)
				}
			})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player kills up to 3 other players.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{killUpTo(3)},
		},
	},
}
