package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Begotten (Basic Monster Card)
//
//	While this is at 1{HP}, it has +1{ATK}.
var begotten = engine.CardDef{
	Ref:     "begotten",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Statics: []engine.Static{{Stat: engine.StatMonsterATK, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, m engine.ObjectID) bool {
		return m == self && g.HP(self) == 1 // "While this is at 1 HP"
	}}},
}
