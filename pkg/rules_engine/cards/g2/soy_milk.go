package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Soy Milk (Passive Treasure Card)
//
//	Monsters have +1{DC} on your turn.
//	Each time you deal combat damage to a monster, deal 1 damage to each other monster.
var soyMilk = engine.CardDef{
	Ref:    "soy_milk",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Statics: []engine.Static{{
		Stat:   engine.StatMonsterDC,
		Amount: 1,
		Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, _ engine.ObjectID) bool {
			return g.Turn.Active == g.Object(self).Controller // "on your turn"
		},
	}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you deal combat damage to a monster, deal 1 damage to each other monster.",
		Trigger: combatDamageToAMonster(0),
		Effects: []engine.Effect{damageOtherMonsters},
	}},
}
