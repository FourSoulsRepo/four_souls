package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum Friend (Active Treasure Card)
//
//	{Tap Effect}Loot 1, then put a loot card from your hand on top of the loot deck.
var bumFriend = engine.CardDef{
	Ref:    "bum_friend",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Loot 1, then put a loot card from your hand on top of the loot deck.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Loot(1), engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					c.G.HandToDeckTop(c.Controller, c.HandCard(a[0]))
				}
			}, engine.HandQuestion("Put which card on top of the loot deck?"))},
		},
	},
}
