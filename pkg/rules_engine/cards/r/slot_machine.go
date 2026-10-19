package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Slot Machine (Room Card)
//
//	{Paid Effect}Pay 3¢: roll-
//	2: Gain 1¢
//	3: Gain 2¢
//	4: Gain 5¢
//	5: Gain 7¢
//	6: Gain 15¢. Put this into discard.
var slotMachine = engine.CardDef{
	Ref:    "slot_machine",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Pay 3¢: Roll- 2: Gain 1¢. 3: Gain 2¢. 4: Gain 5¢. 5: Gain 7¢. 6: Gain 15¢. Put this into discard.",
		Costs: []engine.Cost{engine.PayCents(3)},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(2, 2, engine.GainCents(1)).
			Results(3, 3, engine.GainCents(2)).
			Results(4, 4, engine.GainCents(5)).
			Results(5, 5, engine.GainCents(7)).
			Results(6, 6, engine.GainCents(15), discardThisRoom))},
	}},
}
