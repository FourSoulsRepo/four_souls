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
				{Text: "Put each monster not being attacked on the bottom of the monster deck.", Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
					for _, id := range notAttacked(c.G) {
						c.G.SlotToDeckBottom(id)
					}
					c.G.RefillSlots()
				})}},
				{Text: "Put each shop item on the bottom of the treasure deck.", Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
					for _, id := range c.G.ShopItems() {
						c.G.SlotToDeckBottom(id)
					}
					c.G.RefillSlots()
				})}},
			},
		},
	},
}
