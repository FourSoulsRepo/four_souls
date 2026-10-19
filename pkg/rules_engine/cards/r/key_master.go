package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Key Master (Room Card)
//
//	{Paid Effect}Discard a loot card: Roll. If the discarded card had "key" in its name, roll 3x instead-
//	5-6: Gain +1 treasure.
var keyMaster = engine.CardDef{
	Ref:    "key_master",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "$ Discard a loot card: Roll. If the discarded card had \"key\" in its name, roll 3x instead- 5-6: Gain +1 treasure.",
		Costs:   []engine.Cost{engine.DiscardChosen(0)},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
		Effects: []engine.Effect{rollTimesIfNamed("key", engine.RollTable{}.Results(5, 6, engine.GainTreasure(1)))},
	}},
}
