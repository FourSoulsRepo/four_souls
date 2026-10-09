package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum-Bo! (Passive Treasure Card)
//
//	If you would gain any amount of ¢, this levels up by that much instead.
//	{LV1 Effect}You have +2 to your first attack roll each turn.
//	{LV10 Effect}You have +1{ATK}.
//	{LV25 Effect}You may attack any number of times on your turn.
var bumBo = engine.CardDef{
	Ref:    "bum_bo",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Replacements: []engine.Replacement{{
		Text: "If you would gain any amount of ¢, this levels up by that much instead.",
		When: func(g *engine.Game, self engine.ObjectID, a engine.Action) bool {
			return a.Kind == engine.ActGainCents && a.Amount > 0 && a.Player == g.Object(self).Controller
		},
		Do: func(_ *engine.Game, self engine.ObjectID, a engine.Action) []engine.Action {
			return []engine.Action{{Kind: engine.ActAddCounters, Player: a.Player, Object: self, Amount: a.Amount}}
		},
	}},
	Statics: []engine.Static{{Stat: engine.StatAttackRoll, Amount: 2, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return bumBoLevel(g, self, p, 1) && g.Turn.AttackRolls == 0
	}}, {Stat: engine.StatPlayerATK, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return bumBoLevel(g, self, p, 10)
	}}, {Stat: engine.StatAttacks, Amount: 99, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return bumBoLevel(g, self, p, 25)
	}}},
}

// bumBoLevel: Bum-bo's controller is p and it has at least level counters.
func bumBoLevel(g *engine.Game, self engine.ObjectID, p engine.PlayerID, level int) bool {
	o := g.Object(self)
	return o.Controller == p && o.CountersOf("") >= level
}
