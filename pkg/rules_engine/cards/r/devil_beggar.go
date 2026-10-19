package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Devil Beggar (Room Card)
//
//	{Paid Effect}Pay 1{HP}: Roll-
//	1-2: Loot 2.
//	3-4: Loot 1.
//	5-6: Take 1 damage.
var devilBeggar = engine.CardDef{
	Ref:    "devil_beggar",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Pay 1 HP: Roll- 1-2: Loot 2. 3-4: Loot 1. 5-6: Take 1 damage.",
		Costs: []engine.Cost{engine.PayHPCost(1)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.Loot(2)).
			Results(3, 4, engine.Loot(1)).
			Results(5, 6, youTake(1)))},
	}},
}
