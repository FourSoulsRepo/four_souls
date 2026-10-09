package b2

import engine "github.com/FourSoulsRepo/rules_engine"

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

// orderQuestions: which deck, then which card goes on top, then next,
// and so on. Each answer indexes the cards not placed yet.
func orderQuestions() []engine.Question {
	qs := []engine.Question{engine.DeckQuestion("Look at the top 5 cards of which deck?")}
	for range sleightCards {
		qs = append(qs, engine.Question{
			Text: "Which card goes next, from the top?",
			Options: func(c *engine.Ctx, a []int) []string {
				var out []string
				for _, id := range leftToPlace(c, a) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			},
		})
	}
	return qs
}

// leftToPlace is the looked-at cards not placed by the answers so far.
func leftToPlace(c *engine.Ctx, a []int) []engine.ObjectID {
	left := c.G.DeckTop(c.Deck(a[0]), sleightCards)
	for _, i := range a[1:] {
		if i < 0 {
			break
		}
		left = append(left[:i:i], left[i+1:]...)
	}
	return left
}

func putBack(c *engine.Ctx, a []int) {
	d := c.Deck(a[0])
	left := c.G.DeckTop(d, sleightCards)
	c.G.LookAt(c.Controller, left...)
	var order []engine.ObjectID
	for _, i := range a[1:] {
		if i < 0 {
			break
		}
		order = append(order, left[i])
		left = append(left[:i:i], left[i+1:]...)
	}
	c.G.SetDeckTop(d, order)
}
