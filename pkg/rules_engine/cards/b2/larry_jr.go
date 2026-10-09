package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Larry Jr. (Boss Card)
//
//	While this is at 2{HP} or less, it has +1{DC}.
var larryJr = engine.CardDef{
	Ref:     "larry_jr",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
	Statics: []engine.Static{{Stat: engine.StatMonsterDC, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, m engine.ObjectID) bool {
		return m == self && g.HP(self) <= 2
	}}},
}
