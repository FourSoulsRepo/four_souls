package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Tarot Cloth (Passive Treasure Card)
//
//	Each time a player rolls a ❹, they must give you a loot card.
var tarotCloth = engine.CardDef{
	Ref:    "tarot_cloth",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 4, they must give you a loot card.",
			Takes:   engine.TakesEventPlayer,
			Trigger: engine.OnRollOf(4),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					c.G.GiveHandCard(c.EventPlayer, c.Controller, c.G.Players[c.EventPlayer].Hand[a[0]])
				}
			}, engine.Question{
				Text:   cardkit.AskGiveWhichLootCard,
				Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.EventPlayer },
				Options: func(c *engine.Ctx, _ []int) []engine.Msg {
					if c.EventPlayer == engine.NoPlayer || c.EventPlayer == c.Controller {
						return nil
					}
					var out []engine.Msg
					for _, id := range c.G.Players[c.EventPlayer].Hand {
						out = append(out, engine.CardOption(c.G.Object(id).Card))
					}
					return out
				},
			})},
		},
	},
}
