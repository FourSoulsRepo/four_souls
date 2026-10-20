package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Clog (Boss Card)
//
//	When this dies, each player gives their hand to the player to their left.
var clog = engine.CardDef{
	Ref:     "clog",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, each player gives their hand to the player to their left.",
		Takes:   engine.TakesOthers,
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			hands := make([][]engine.ObjectID, len(c.G.Players))
			for i, pl := range c.G.Players {
				hands[i] = append([]engine.ObjectID(nil), pl.Hand...)
			}
			for i, hand := range hands {
				for _, id := range hand {
					c.G.GiveHandCard(engine.PlayerID(i), c.G.PlayerAfter(engine.PlayerID(i)), id)
				}
			}
		})},
	}},
}
