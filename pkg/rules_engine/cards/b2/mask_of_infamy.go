package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mask Of Infamy (Boss Card)
//
//	While this is at 1{HP}, it has +2{DC}.
var maskOfInfamy = engine.CardDef{
	Ref:     "mask_of_infamy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Statics: []engine.Static{{Stat: engine.StatMonsterDC, Amount: 2, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, m engine.ObjectID) bool {
		return m == self && g.HP(self) == 1
	}}},
}
