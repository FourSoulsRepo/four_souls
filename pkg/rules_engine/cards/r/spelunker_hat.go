package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Spelunker Hat (Passive Treasure Card)
//
//	At the start of your turn, you may look at the top 3 cards of a deck and put them back in any order.
var spelunkerHat = engine.CardDef{
	Ref:    "spelunker_hat",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn, you may look at the top 3 cards of a deck and put them back in any order.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if deckChosen(c, a[0]) && a[1] == 0 {
				c.G.SetDeckTop(c.Deck(a[0]), c.Order)
			}
		}, mayDeckQuestion("Look at the top 3 cards of which deck?"), engine.Question{
			Text: cardkit.PutBackText,
			Order: func(c *engine.Ctx, a []int) []engine.ObjectID {
				if !deckChosen(c, a[0]) {
					return nil
				}
				return c.G.DeckTop(c.Deck(a[0]), 3)
			},
		})},
	}},
}
