package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Battery (Active Treasure Card)
//
//	{Tap Effect}Recharge another item.
var theBattery = engine.CardDef{
	Ref:    "the_battery",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Recharge another item.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, otherActiveItem)},
			Effects: []engine.Effect{engine.RechargeTarget(0)},
		},
	},
}

// otherActiveItem: "another item" with a ↷ ability (19.95).
func otherActiveItem(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
	return engine.NotThis(g, self, c) && engine.ActiveItem(g, self, c)
}
