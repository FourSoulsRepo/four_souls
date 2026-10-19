package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Backstabber (Passive Treasure Card)
//
//	You may attack players who control more souls than you. They have 3+{DC} for the attack.
//	Each time you kill a player, loot 1.
var backstabber = engine.CardDef{
	Ref:    "backstabber",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Statics: []engine.Static{{Stat: engine.StatPlayerAttackDC, Amount: 3, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		me := g.Object(self).Controller
		return me == g.Turn.Active && g.SoulValue(p) > g.SoulValue(me)
	}}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you kill a player, loot 1.",
		Trigger: whenYouKillAPlayer(),
		Effects: []engine.Effect{engine.Loot(1)},
	}},
}
