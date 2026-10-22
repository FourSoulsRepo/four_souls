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

func flushMonsters(c *engine.Ctx) []engine.ObjectID { return notAttacked(c.G) }

func flushShop(c *engine.Ctx) []engine.ObjectID { return c.G.ShopItems() }

// flushTo puts the cards on the bottom of the deck in random order
// (Q-22); then the slots refill.
func flushTo(d engine.DeckKind, cards func(c *engine.Ctx) []engine.ObjectID) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.ToBottomRandom(d, cards(c))
		c.G.RefillSlots()
	})
}
