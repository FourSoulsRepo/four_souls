package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bob’s Brain (Passive Treasure Card)
//
//	Each time you declare an attack, roll-
//	1-2: Deal 1 damage to a monster.
//	3-4: Deal 1 damage to a player.
//	5-6: Take 1 damage.
var bobsBrain = engine.CardDef{
	Ref:    "bobs_brain",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you declare an attack, roll- 1-2: Deal 1 damage to a monster. 3-4: Deal 1 damage to a player. 5-6: Take 1 damage.",
			Trigger: whenYouDeclareAnAttack(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
				Results(1, 2, damageAMonster(1)).
				Results(3, 4, damageAPlayer(1, false)).
				Results(5, 6, engine.DealDamage(1, engine.You)))},
		},
	},
}

func whenYouDeclareAnAttack() engine.Trigger {
	return engine.Trigger{On: engine.EvAttackDeclared, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Player == g.Object(self).Controller
	}}
}
