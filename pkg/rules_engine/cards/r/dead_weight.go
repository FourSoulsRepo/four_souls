package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Dead Weight (Eternal Treasure Card)
//
//	You take no combat damage on attack rolls of 1.
//	{Tap Effect}Cancel the triggered ability of a monster or non-eternal item.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var deadWeight = engine.CardDef{
	Ref:       "dead_weight",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	Outside:   true,
	Tap:       true,
	CombatMod: noCombatDamageTakenOn(1),
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Cancel the triggered ability of a monster or non-eternal item.",
			Costs: []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetStackTrigger, func(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
				it, ok := g.StackItemByID(c.StackID)
				if !ok {
					return false
				}
				d := g.Def(it.Source)
				return d.Kind == engine.MonsterCard || (d.Kind == engine.TreasureCard && !g.Eternal(it.Source))
			})},
			Effects: []engine.Effect{engine.CancelTarget(0)},
		},
		cardkit.RechargeAtEndOfTurn,
	},
}
