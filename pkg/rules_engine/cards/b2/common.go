package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pieces many Base Game cards share.

// extraLoot is every Base Game character's ↷ ability. It can be used on
// any player's turn, in response to anything.
var extraLoot = engine.Ability{
	Kind:    engine.Activated,
	Text:    "↷: Play an additional loot card this turn.",
	Costs:   []engine.Cost{engine.Tap()},
	Effects: []engine.Effect{engine.AddLootPlays(1)},
}

// rechargeAtEndOfTurn: "At the end of your turn, recharge this".
var rechargeAtEndOfTurn = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "At the end of your turn, recharge this.",
	Trigger: engine.AtEndOfYourTurn(),
	Effects: []engine.Effect{engine.RechargeSelf()},
}

// discardOne is "Discard 1 loot card." for the controller.
var discardOne = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.DiscardFromHand(c.Controller, c.HandCard(a[0]))
	}
}, engine.HandQuestion("Discard which loot card?"))

// damageEqualToRoll deals damage equal to the roll to target 0.
func damageEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.DealDamage(r, 0))
	}
	return t
}

// endYourTurn: "If it's your turn, cancel everything that hasn't
// resolved and end it".
var endYourTurn = engine.EffectFunc(func(c *engine.Ctx) {
	if c.Controller == c.G.Turn.Active {
		c.G.EndTurnNow()
	}
})

// lookMayBottom: "look at the top card of the deck. You may put it on
// the bottom".
func lookMayBottom(d engine.DeckKind) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		top := c.G.DeckTop(d, 1)
		if len(top) == 0 {
			return
		}
		c.G.LookAt(c.Controller, top...)
		if a[0] == 1 {
			c.G.DeckToBottom(d, top[0])
		}
	}, engine.Question{
		Text: "Put the top card on the bottom?",
		Options: func(c *engine.Ctx, _ []int) []string {
			top := c.G.DeckTop(d, 1)
			if len(top) == 0 {
				return nil
			}
			card := string(c.G.Object(top[0]).Card)
			return []string{"keep " + card + " on top", "put " + card + " on the bottom"}
		},
	})
}

// oneOnTop: "Look at the top 5 cards of the deck. Put 1 on top and the
// rest on the bottom".
func oneOnTop(d engine.DeckKind) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		top := c.G.DeckTop(d, 5)
		c.G.LookAt(c.Controller, top...)
		if a[0] < 0 {
			return
		}
		for i, id := range top {
			if i != a[0] {
				c.G.DeckToBottom(d, id)
			}
		}
	}, engine.Question{
		Text: "Which card stays on top? The rest go to the bottom.",
		Options: func(c *engine.Ctx, _ []int) []string {
			var out []string
			for _, id := range c.G.DeckTop(d, 5) {
				out = append(out, string(c.G.Object(id).Card))
			}
			return out
		},
	})
}
