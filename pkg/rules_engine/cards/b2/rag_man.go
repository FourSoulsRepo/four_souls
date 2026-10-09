package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rag Man (Boss Card)
//
//	When this dies, after gaining rewards, the active player rolls-
//	1 or 6: Put this on top of the monster deck.
var ragMan = engine.CardDef{
	Ref:     "rag_man",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      2,
	DC:      3,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, after gaining rewards, the active player rolls- 1 or 6: Put this on top of the monster deck.",
			Trigger: engine.AfterThisRewards(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 1, ragManBack).Results(6, 6, ragManBack))},
		},
	},
}

// ragManBack puts the dead Rag Man on top of the monster deck.
var ragManBack = engine.EffectFunc(func(c *engine.Ctx) {
	if c.G.Object(c.EventObject).Zone.Kind == engine.ZoneOutside {
		c.G.PutIntoDeck(engine.MonsterDeck, c.EventObject, 0)
	}
})
