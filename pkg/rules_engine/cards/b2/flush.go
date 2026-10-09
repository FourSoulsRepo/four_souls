package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Flush! (Active Treasure Card)
//
//	{Tap Effect}Choose one- Put each monster not being attacked on the bottom of the monster deck. Put each shop item on the bottom of the treasure deck.
var flush = engine.CardDef{
	Ref:    "flush",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Choose one- Put each monster not being attacked on the bottom of the monster deck. Put each shop item on the bottom of the treasure deck.",
			Costs: []engine.Cost{engine.Tap()},
			Modes: []engine.Mode{
				{Text: "Put each monster not being attacked on the bottom of the monster deck.", Effects: []engine.Effect{flushTo(engine.MonsterDeck, flushMonsters)}},
				{Text: "Put each shop item on the bottom of the treasure deck.", Effects: []engine.Effect{flushTo(engine.TreasureDeck, flushShop)}},
			},
		},
	},
}

func flushMonsters(c *engine.Ctx, _ []int) []engine.ObjectID { return notAttacked(c.G) }

func flushShop(c *engine.Ctx, _ []int) []engine.ObjectID { return c.G.ShopItems() }

// flushTo puts the cards on the bottom of the deck in the order the
// player picks, the first one lowest; then the slots refill.
func flushTo(d engine.DeckKind, cards func(c *engine.Ctx, a []int) []engine.ObjectID) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		toBottom(c.G, d, engine.Ordered(cards(c, nil), a))
		c.G.RefillSlots()
	}, engine.OrderQuestions("Which card goes lowest?", "the rest in slot order", 0, 6, nil, cards)...)
}
