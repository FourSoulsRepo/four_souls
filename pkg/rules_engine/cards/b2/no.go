package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// No! (Active Treasure Card)
//
//	{Tap Effect}Cancel the ↷ or $ ability of an item.
var no = engine.CardDef{
	Ref:    "no",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Cancel the ↷ or $ ability of an item.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetStackAbility, itemAbility)},
			Effects: []engine.Effect{engine.CancelTarget(0)},
		},
	},
}

// itemAbility keeps activated abilities of items, not loot cards.
func itemAbility(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	it, ok := g.StackItemByID(c.StackID)
	return ok && it.Kind == engine.StackAbility
}
