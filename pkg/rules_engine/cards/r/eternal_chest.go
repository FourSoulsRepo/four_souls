package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Eternal Chest (Room Card)
//
//	{Paid Effect}Discard a loot card: Roll-
//	3: Gain 4¢.
//	4-5: Gain 7¢.
//	6: Gain +1 treasure and loot 1.
var eternalChest = engine.CardDef{
	Ref:    "eternal_chest",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "$ Discard a loot card: Roll- 3: Gain 4¢. 4-5: Gain 7¢. 6: Gain +1 treasure and loot 1.",
		Costs:   []engine.Cost{engine.DiscardChosen(0)},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(3, 3, engine.GainCents(4)).
			Results(4, 5, engine.GainCents(7)).
			Results(6, 6, engine.GainTreasure(1), engine.Loot(1)))},
	}},
}
