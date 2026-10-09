package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Curse (Eternal Treasure Card)
//
//	At the start of your turn, put the top card of a deck into discard.
//	{Tap Effect}Put the top card of any discard on top of its deck.
//	-Eternal- This can't be destroyed or put into discard.
var theCurse = engine.CardDef{
	Ref:     "the_curse",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, put the top card of a deck into discard.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					c.G.MillTop(c.Deck(a[0]))
				}
			}, engine.DeckQuestion("Put the top card of which deck into discard?"))},
		},
		{
			Kind:    engine.Activated,
			Text:    "↷: Put the top card of any discard on top of its deck.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Ask(discardBack, discardQuestion)},
		},
	},
}

// discardQuestion offers every discard pile that has cards.
var discardQuestion = engine.Question{
	Text: "Put the top card of which discard on top of its deck?",
	Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, d := range piles(c) {
			out = append(out, d.String()+" discard")
		}
		return out
	},
}

func piles(c *engine.Ctx) []engine.DeckKind {
	var out []engine.DeckKind
	for _, d := range c.G.DecksInPlay() {
		if len(c.G.Discards[d]) > 0 {
			out = append(out, d)
		}
	}
	return out
}

func discardBack(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.DiscardTopToDeck(piles(c)[a[0]])
	}
}
