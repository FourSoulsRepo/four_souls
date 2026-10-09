package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gemini (Boss Card)
//
//	While this is at 1{HP}, it has +1{ATK}.
var gemini = engine.CardDef{
	Ref:     "gemini",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Statics: []engine.Static{{Stat: engine.StatMonsterATK, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, m engine.ObjectID) bool {
		return m == self && g.HP(self) == 1
	}}},
}
