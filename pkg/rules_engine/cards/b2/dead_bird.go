package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	deadBirdDont = engine.NewKey("card.dead_bird.dont")
)

// Dead Bird (Passive Treasure Card)
//
//	Each time a player rolls a ❸, you may look at their hand and steal a loot card from them.
var deadBird = engine.CardDef{
	Ref:    "dead_bird",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 3, you may look at their hand and steal a loot card from them.",
			Takes:   engine.TakesEventPlayer,
			Trigger: engine.OnRollOf(3),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				from := c.EventPlayer
				if from == engine.NoPlayer || from == c.Controller || a[0] < 0 || a[0] >= len(c.G.Players[from].Hand) {
					return
				}
				c.G.LookAt(c.Controller, c.G.Players[from].Hand...)
				c.G.GiveHandCard(from, c.Controller, c.G.Players[from].Hand[a[0]])
			}, engine.Question{Text: cardkit.AskStealWhichCard, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				if c.EventPlayer == engine.NoPlayer || c.EventPlayer == c.Controller {
					return nil
				}
				var out []engine.Msg
				for _, id := range c.G.Players[c.EventPlayer].Hand {
					out = append(out, engine.CardOption(c.G.Object(id).Card))
				}
				return append(out, deadBirdDont.Msg())
			}})},
		},
	},
}
