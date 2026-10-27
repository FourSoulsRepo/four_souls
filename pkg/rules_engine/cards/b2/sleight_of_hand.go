package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	sleightOfHandLookAtThe = engine.NewKey("card.sleight_of_hand.look_at_the_top_5")
)

// Sleight of Hand (Eternal Treasure Card)
//
//	{Tap Effect}Look at the top 5 cards of a deck. Put them back in any order.
//	-Eternal- This can't be destroyed or put into discard.
var sleightOfHand = engine.CardDef{
	Ref:     "sleight_of_hand",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Look at the top 5 cards of a deck. Put them back in any order.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Ask(putBack, orderQuestions()...)},
		},
	},
}

const sleightCards = 5

// orderQuestions: which deck, then the order of its top cards.
func orderQuestions() []engine.Question {
	return []engine.Question{
		engine.DeckQuestion(sleightOfHandLookAtThe),
		{
			Text:  cardkit.AskPutBack,
			Order: func(c *engine.Ctx, a []int) []engine.ObjectID { return c.G.DeckTop(c.Deck(a[0]), sleightCards) },
		},
	}
}

func putBack(c *engine.Ctx, a []int) {
	c.G.SetDeckTop(c.Deck(a[0]), c.Order)
}
