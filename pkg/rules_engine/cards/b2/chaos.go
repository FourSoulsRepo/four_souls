package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Chaos (Active Treasure Card)
//
//	{Tap Effect}Each player gives their hand to the player to their left.
var chaos = engine.CardDef{
	Ref:    "chaos",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Each player gives their hand to the player to their left.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				hands := make([][]engine.ObjectID, len(c.G.Players))
				for i, pl := range c.G.Players {
					hands[i] = append([]engine.ObjectID(nil), pl.Hand...)
				}
				for i, hand := range hands {
					to := engine.PlayerID((i + 1) % len(hands))
					for _, id := range hand {
						c.G.GiveHandCard(engine.PlayerID(i), to, id)
					}
				}
			})},
		},
	},
}
