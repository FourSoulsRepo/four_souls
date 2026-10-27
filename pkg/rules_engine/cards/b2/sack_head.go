package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	sackHeadPutItOn = engine.NewKey("card.sack_head.put_it_on_the_bottom")
)

// Sack Head (Active Treasure Card)
//
//	{Tap Effect}Look at the top card of a deck. You may put that card on the bottom of that deck.
var sackHead = engine.CardDef{
	Ref:    "sack_head",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Look at the top card of a deck. You may put that card on the bottom of that deck.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] < 0 {
					return
				}
				d := c.Deck(a[0])
				top := c.G.DeckTop(d, 1)
				if len(top) == 0 {
					return
				}
				c.G.LookAt(c.Controller, top...)
				if a[1] == 1 {
					c.G.DeckToBottom(d, top[0])
				}
			}, engine.DeckQuestion(cardkit.AskLookTop), engine.Question{
				Text: sackHeadPutItOn,
				Subject: func(c *engine.Ctx, a []int) engine.ObjectID {
					return c.G.DeckTop(c.Deck(a[0]), 1)[0]
				},
				Options: func(c *engine.Ctx, a []int) []engine.Msg {
					top := c.G.DeckTop(c.Deck(a[0]), 1)
					if len(top) == 0 {
						return nil
					}
					return cardkit.KeepOrBottom(c.G.Object(top[0]).Card)
				},
			})},
		},
	},
}
