package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Incubus (Eternal Treasure Card)
//
//	{Tap Effect}Choose one- Look at a player's hand. You may swap a card
//	from your hand with one of theirs. Loot 1, then put a card from your
//	hand on top of the loot deck.
//	-Eternal- This can't be destroyed or put into discard.
var incubus = engine.CardDef{
	Ref:     "incubus",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Choose one- Look at a player's hand. You may swap a card from your hand with one of theirs. Loot 1, then put a card from your hand on top of the loot deck.",
			Costs: []engine.Cost{engine.Tap()},
			Modes: []engine.Mode{
				{
					Text:    "Look at a player's hand. You may swap a card from your hand with one of theirs.",
					Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
					Effects: []engine.Effect{engine.Ask(swapHands, theirCard, yourCard)},
				},
				{
					Text: "Loot 1, then put a card from your hand on top of the loot deck.",
					Effects: []engine.Effect{
						engine.Loot(1),
						engine.Ask(func(c *engine.Ctx, a []int) {
							if a[0] >= 0 {
								c.G.HandToDeckTop(c.Controller, c.HandCard(a[0]))
							}
						}, engine.HandQuestion("Put which card on top of the loot deck?")),
					},
				},
			},
		},
	},
}

// theirCard shows the other player's hand: pick a card to take, or keep
// your cards. The last option is always "don't swap".
var theirCard = engine.Question{
	Text: "Their hand. Swap for which card?",
	Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range c.G.Players[c.Targets[0].Player].Hand {
			out = append(out, string(c.G.Object(id).Card))
		}
		return append(out, "don't swap")
	},
}

// yourCard asks which of your cards to give; skipped when not swapping.
var yourCard = engine.Question{
	Text: "Give which of your cards?",
	Options: func(c *engine.Ctx, a []int) []string {
		if a[0] == len(c.G.Players[c.Targets[0].Player].Hand) {
			return nil // "don't swap"
		}
		return engine.HandQuestion("").Options(c, a)
	},
}

func swapHands(c *engine.Ctx, a []int) {
	them := c.Targets[0].Player
	c.G.LookAt(c.Controller, c.G.Players[them].Hand...)
	if a[1] < 0 {
		return
	}
	theirs, yours := c.G.Players[them].Hand[a[0]], c.HandCard(a[1])
	// A swap: both cards change hands at the same time (R-MECH-39).
	c.G.GiveHandCard(them, c.Controller, theirs)
	c.G.GiveHandCard(c.Controller, them, yours)
}
