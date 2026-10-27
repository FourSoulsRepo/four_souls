package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	incubusGiveWhichOf   = engine.NewKey("card.incubus.give_which_of_your_cards")
	incubusTheirHandSwap = engine.NewKey("card.incubus.their_hand_swap_for_which")
)

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
						}, engine.HandQuestion(cardkit.AskPutWhichCardOnTop)),
					},
				},
			},
		},
	},
}

// theirCard shows the other player's hand: pick a card to take, or keep
// your cards. The last option is always "don't swap".
var theirCard = engine.Question{
	Text: incubusTheirHandSwap,
	Options: func(c *engine.Ctx, _ []int) []engine.Msg {
		var out []engine.Msg
		for _, id := range c.G.Players[c.Targets[0].Player].Hand {
			out = append(out, engine.CardOption(c.G.Object(id).Card))
		}
		return append(out, cardkit.AskDontSwap.Msg())
	},
}

// yourCard asks which of your cards to give; skipped when not swapping.
var yourCard = engine.Question{
	Text: incubusGiveWhichOf,
	Options: func(c *engine.Ctx, a []int) []engine.Msg {
		if a[0] == len(c.G.Players[c.Targets[0].Player].Hand) {
			return nil // "don't swap"
		}
		return engine.HandQuestion(incubusGiveWhichOf).Options(c, a)
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
