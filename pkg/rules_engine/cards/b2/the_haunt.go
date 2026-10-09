package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Haunt (Boss Card)
//
//	Every other time this takes damage each turn, it gains +1{DC} till end of turn.
var theHaunt = engine.CardDef{
	Ref:     "the_haunt",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Every other time this takes damage each turn, it gains +1 DC till end of turn.",
			Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Object == self && g.Object(self).HitsThisTurn%2 == 0
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatMonsterDC, Player: engine.NoPlayer, Object: c.Source, Amount: 1})
			})},
		},
	},
}
