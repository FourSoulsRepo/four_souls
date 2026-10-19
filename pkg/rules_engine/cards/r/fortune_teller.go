package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Fortune Teller (Room Card)
//
//	{Paid Effect}Pay 3¢: roll-
//	3: Loot 1.
//	4: Loot 2.
//	5: Loot 3.
//	6: Loot 6. Put this into discard.
var fortuneTeller = engine.CardDef{
	Ref:    "fortune_teller",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Pay 3¢: Roll- 3: Loot 1. 4: Loot 2. 5: Loot 3. 6: Loot 6. Put this into discard.",
		Costs: []engine.Cost{engine.PayCents(3)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(3, 3, engine.Loot(1)).
			Results(4, 4, engine.Loot(2)).
			Results(5, 5, engine.Loot(3)).
			Results(6, 6, engine.Loot(6), discardThisRoom))},
	}},
}
