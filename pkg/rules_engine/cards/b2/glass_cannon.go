package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Glass Cannon (Active Treasure Card)
//
//	{Tap Effect}Destroy another item, then roll-
//	1-5: Destroy this and loot 2.
//	6: Recharge this.
var glassCannon = engine.CardDef{
	Ref:    "glass_cannon",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Destroy another item, then roll- 1-5: Destroy this and loot 2. 6: Recharge this.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, engine.NotThis)},
			Effects: []engine.Effect{engine.DestroyTarget(0), engine.Roll(engine.RollTable{}.
				Results(1, 5, engine.DestroyThis(), engine.Loot(2)).
				Results(6, 6, engine.RechargeSelf()))},
		},
	},
}
