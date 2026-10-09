package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Delirium (Epic Boss Card)
//
//	Other monsters have +1{DC}.
//	When this dies, put it in the monster deck 6 cards from the top.
var delirium = engine.CardDef{
	Ref:     "delirium",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     3,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
	Statics: []engine.Static{{Stat: engine.StatMonsterDC, Amount: 1, Applies: func(_ *engine.Game, self engine.ObjectID, _ engine.PlayerID, m engine.ObjectID) bool {
		return m != self
	}}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, put it in the monster deck 6 cards from the top.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if c.G.Object(c.EventObject).Zone.Kind == engine.ZoneOutside {
					c.G.PutIntoDeck(engine.MonsterDeck, c.EventObject, 5)
				}
			})},
		},
	},
}
