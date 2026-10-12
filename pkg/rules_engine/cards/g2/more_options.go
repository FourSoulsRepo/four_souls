package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// More Options (Passive Treasure Card)
//
//	At the start of your turn, you may look at the top 2 cards of the treasure deck.
//	You may purchase an additional time on your turn.
var moreOptions = engine.CardDef{
	Ref:     "more_options",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPurchases, 1)},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn, you may look at the top 2 cards of the treasure deck.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.LookAt(c.Controller, c.G.DeckTop(engine.TreasureDeck, 2)...)
		})},
	}},
}
