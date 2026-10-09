package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curved Horn (Trinket Card)
//
//	Gain +1{ATK} for your first attack roll each turn.
//	-Trinket- This loot becomes an item under your control when it resolves.
var curvedHorn = engine.CardDef{
	Ref:     "curved_horn",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Statics: []engine.Static{{
		Stat:   engine.StatPlayerATK,
		Amount: 1,
		// Only while the first attack roll of the turn deals its damage.
		Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
			return g.Object(self).Controller == p && p == g.Turn.Active && g.Attack.On && g.Turn.AttackRolls == 1
		},
	}},
	Abilities: []engine.Ability{},
}
