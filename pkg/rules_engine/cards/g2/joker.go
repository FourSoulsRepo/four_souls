package g2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Joker (Wildcard Card)
//
//	Look at a player's hand. You may steal a loot card from them.
var joker = engine.CardDef{
	Ref:    "joker",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.LootAbility,
		Text:  "Look at a player's hand. You may steal a loot card from them.",
		Takes: engine.TakesOthers,
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			ps := livingPlayers(c.G, c.Controller)
			if a[0] < 0 || a[0] >= len(ps) {
				return
			}
			hand := c.G.Players[ps[a[0]]].Hand
			c.G.LookAt(c.Controller, hand...)
			if a[1] >= 0 && a[1] < len(hand) {
				c.G.GiveHandCard(ps[a[0]], c.Controller, hand[a[1]])
			}
		}, pickAnotherPlayer(cardkit.AskLookAtWhoseHand), engine.Question{Text: cardkit.AskStealWhichCard, Options: func(c *engine.Ctx, a []int) []engine.Msg {
			if a[0] < 0 {
				return nil
			}
			var out []engine.Msg
			for _, id := range c.G.Players[livingPlayers(c.G, c.Controller)[a[0]]].Hand {
				out = append(out, engine.CardOption(c.G.Object(id).Card))
			}
			return append(out, engine.AskNone.Msg())
		}})},
	}},
}
