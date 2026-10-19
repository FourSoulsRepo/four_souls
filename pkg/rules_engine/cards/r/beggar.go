package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Beggar (Room Card)
//
//	{Paid Effect}Pay 3¢: Roll-
//	3-4: Gain 2¢
//	5: Loot 2.
//	6: Gain +1 treasure. Put this into discard.
var beggar = engine.CardDef{
	Ref:    "beggar",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Pay 3¢: Roll- 3-4: Gain 2¢. 5: Loot 2. 6: Gain +1 treasure. Put this into discard.",
		Costs: []engine.Cost{engine.PayCents(3)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(3, 4, engine.GainCents(2)).
			Results(5, 5, engine.Loot(2)).
			Results(6, 6, engine.GainTreasure(1), discardThisRoom))},
	}},
}
