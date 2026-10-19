package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Donation (Room Card)
//
//	{Paid Effect}Pay 1{HP}:
//	Gain 3¢.
var bloodDonation = engine.CardDef{
	Ref:    "blood_donation",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "$ Pay 1 HP: Gain 3¢.",
		Costs:   []engine.Cost{engine.PayHPCost(1)},
		Effects: []engine.Effect{engine.GainCents(3)},
	}},
}
