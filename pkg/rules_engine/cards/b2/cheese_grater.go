package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	cheeseGraterDiscard      = engine.NewKey("card.cheese_grater.discard")
	cheeseGraterPutBack      = engine.NewKey("card.cheese_grater.put_back")
	cheeseGraterPutItBack    = engine.NewKey("card.cheese_grater.put_it_back_or_into")
	cheeseGraterRevealTheTop = engine.NewKey("card.cheese_grater.reveal_the_top_card_of")
)

// Cheese Grater (Passive Treasure Card)
//
//	Each time a player rolls a ❻, reveal the top card of any deck. Put it back or put it into discard.
var cheeseGrater = engine.CardDef{
	Ref:    "cheese_grater",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 6, reveal the top card of any deck. Put it back or put it into discard.",
			Trigger: engine.OnRollOf(6),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] < 0 {
					return
				}
				d := c.Deck(a[0])
				c.G.RevealTop(d)
				if a[1] == 1 {
					c.G.MillTop(d)
				}
			}, engine.DeckQuestion(cheeseGraterRevealTheTop), engine.Question{Text: cheeseGraterPutItBack, Options: func(c *engine.Ctx, a []int) []engine.Msg {
				if a[0] < 0 || len(c.G.DeckTop(c.Deck(a[0]), 1)) == 0 {
					return nil
				}
				return engine.Options(cheeseGraterPutBack, cheeseGraterDiscard)
			}})},
		},
	},
}
