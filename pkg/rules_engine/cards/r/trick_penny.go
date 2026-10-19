package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Trick Penny (Passive Treasure Card)
//
//	Shop items you purchase cost 1¢ less.
//	$ items you control cost 1¢ less to activate.
//	Each time you purchase from the shop or treasure deck, gain 1¢.
//	At the end of your turn, gain 1¢.
var trickPenny = engine.CardDef{
	Ref:    "trick_penny",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Statics: []engine.Static{
		engine.YouHave(engine.StatShopPrice, -1),
		engine.YouHave(engine.StatPaidDiscount, 1),
	},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time you purchase from the shop or treasure deck, gain 1¢.",
			Trigger: engine.Trigger{On: engine.EvPurchased, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Player == g.Object(self).Controller
			}},
			Effects: []engine.Effect{engine.GainCents(1)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, gain 1¢.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.GainCents(1)},
		},
	},
}
