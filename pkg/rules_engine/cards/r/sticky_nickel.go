package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sticky Nickel (Nickel Card)
//
//	Roll-
//	1-3: Give this to the player to your left.
//	4-6: Gain 5¢.
var stickyNickel = engine.CardDef{
	Ref:    "sticky_nickel",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Roll- 1-3: Give this to the player to your left. 4-6: Gain 5¢.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 3, engine.EffectFunc(func(c *engine.Ctx) {
				// The card went to the loot discard when it resolved; it goes
				// from there to the left player's hand (S-YUG).
				pile := c.G.Discards[engine.LootDeck]
				for i := len(pile) - 1; i >= 0; i-- {
					if c.G.Object(pile[i]).Card == "sticky_nickel" {
						c.G.DiscardToHand(engine.LootDeck, pile[i], c.G.PlayerAfter(c.Controller))
						return
					}
				}
			})).
			Results(4, 6, engine.GainCents(5)))},
	}},
}
