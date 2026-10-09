package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Donation Machine (Paid Treasure Card)
//
//	{Paid Effect}Give another non-eternal item you control to another player:
//	Gain 8¢.
var donationMachine = engine.CardDef{
	Ref:    "donation_machine",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "Give another non-eternal item you control to another player: Gain 8¢.",
			Costs:   []engine.Cost{engine.GiveChosen(0, 1)},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetYourItem, otherNonEternal), engine.Choose(engine.TargetOtherPlayer)},
			Effects: []engine.Effect{engine.GainCents(8)},
		},
	},
}

// otherNonEternal keeps your non-eternal items other than this one.
func otherNonEternal(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
	return c.Object != self && !g.Eternal(c.Object)
}
