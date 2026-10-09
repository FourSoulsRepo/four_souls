package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Portable Slot Machine (Paid Treasure Card)
//
//	{Paid Effect}Pay 3¢: Roll-
//	1-2: Loot 1.
//	3-4: Gain 4¢.
var portableSlotMachine = engine.CardDef{
	Ref:    "portable_slot_machine",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Pay 3¢: Roll- 1-2: Loot 1. 3-4: Gain 4¢.",
			Costs:   []engine.Cost{engine.PayCents(3)},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 2, engine.Loot(1)).Results(3, 4, engine.GainCents(4)))},
		},
	},
}
