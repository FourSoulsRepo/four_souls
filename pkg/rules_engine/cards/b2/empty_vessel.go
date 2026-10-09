package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Empty Vessel (Passive Treasure Card)
//
//	When you have 0 loot cards in your hand, you have +1{ATK}.
//	While you have 0¢, you have +1 to your attack rolls.
var emptyVessel = engine.CardDef{
	Ref:    "empty_vessel",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Statics: []engine.Static{{Stat: engine.StatPlayerATK, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return g.Object(self).Controller == p && len(g.Players[p].Hand) == 0
	}}, {Stat: engine.StatAttackRoll, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return g.Object(self).Controller == p && g.Players[p].Cents == 0
	}}},
}
