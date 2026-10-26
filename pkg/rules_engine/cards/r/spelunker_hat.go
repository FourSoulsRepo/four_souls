package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	spelunkerHatLookAtThe = engine.NewKey("card.spelunker_hat.look_at_the_top_3")
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
		}, mayDeckQuestion(spelunkerHatLookAtThe), engine.Question{
			Text: cardkit.AskPutBack,
			Order: func(c *engine.Ctx, a []int) []engine.ObjectID {
				if !deckChosen(c, a[0]) {
					return nil
				}
				return c.G.DeckTop(c.Deck(a[0]), 3)
			},
		})},
	}},
}
