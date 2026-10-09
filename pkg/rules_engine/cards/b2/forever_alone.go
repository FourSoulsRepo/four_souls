package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Forever Alone (Eternal Treasure Card)
//
//	{Tap Effect}Choose one- Steal 1¢ from another player. Look at the top
//	card of a deck. Discard a loot card, then loot 1.
//	Each time you take damage, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var foreverAlone = engine.CardDef{
	Ref:     "forever_alone",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Choose one- Steal 1¢ from another player. Look at the top card of a deck. Discard a loot card, then loot 1.",
			Costs: []engine.Cost{engine.Tap()},
			Modes: []engine.Mode{
				{
					Text:    "Steal 1¢ from another player.",
					Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
					Effects: []engine.Effect{engine.StealCents(1, 0)},
				},
				{
					Text: "Look at the top card of a deck.",
					Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
						if a[0] >= 0 {
							c.G.LookAt(c.Controller, c.G.DeckTop(c.Deck(a[0]), 1)...)
						}
					}, engine.DeckQuestion("Look at the top card of which deck?"))},
				},
				{
					Text: "Discard a loot card, then loot 1.",
					Effects: []engine.Effect{
						engine.Ask(func(c *engine.Ctx, a []int) {
							if a[0] >= 0 {
								c.G.DiscardFromHand(c.Controller, c.HandCard(a[0]))
							}
						}, engine.HandQuestion("Discard which loot card?")),
						engine.Loot(1),
					},
				},
			},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, recharge this.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{engine.RechargeSelf()},
		},
	},
}
