package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gurglings (Boss Card)
//
//	While this is at 2{HP} or less, it has -1{DC}.
var gurglings = engine.CardDef{
	Ref:     "gurglings",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Statics: []engine.Static{{Stat: engine.StatMonsterDC, Amount: -1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, m engine.ObjectID) bool {
		return m == self && g.HP(self) <= 2 // "While this is at 2 HP or less"
	}}},
}
