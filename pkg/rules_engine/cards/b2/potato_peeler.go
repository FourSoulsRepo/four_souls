package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Potato Peeler (Active Treasure Card)
//
//	{Tap Effect}Put the top card of each deck into discard.
var potatoPeeler = engine.CardDef{
	Ref:    "potato_peeler",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Put the top card of each deck into discard.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				for _, d := range c.G.DecksInPlay() {
					c.G.MillTop(d)
				}
			})},
		},
	},
}
