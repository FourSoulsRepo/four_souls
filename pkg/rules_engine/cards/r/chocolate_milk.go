package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Chocolate Milk (Passive Treasure Card)
//
//	At the start of your turn, put a counter on this. Then if this has 3+ counters, remove all counters from this and deal 1 damage to a player or monster.
//	While this has 2 counters on it, you have +1{ATK}.
var chocolateMilk = engine.CardDef{
	Ref:    "chocolate_milk",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn, put a counter on this. Then if this has 3+ counters, remove all counters from this and deal 1 damage to a player or monster.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{engine.AddCounters(1), engine.EffectFunc(func(c *engine.Ctx) {
			if n := c.G.Object(c.Source).CountersOf(""); n >= 3 {
				c.G.AddCountersTo(c.Source, -n)
				c.Do(damageAnything(1))
			}
		})},
	}},
	Statics: []engine.Static{{Stat: engine.StatPlayerATK, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		o := g.Object(self)
		return o.Controller == p && o.CountersOf("") == 2
	}}},
}
