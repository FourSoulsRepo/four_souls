package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	fetalHauntGiveItTo = engine.NewKey("card.fetal_haunt.give_it_to_which_player")
)

// Fetal Haunt (Passive Treasure Card)
//
//	{Curse Effect}At the end of your turn, give another player a loot card.
//	When you die, before paying penalties, give this to another player.
var fetalHaunt = engine.CardDef{
	Ref:    "fetal_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, give another player a loot card.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				hand, others := c.G.Players[c.Controller].Hand, otherPlayers(c)
				if a[0] < 0 || a[0] >= len(hand) || a[1] < 0 || a[1] >= len(others) {
					return
				}
				c.G.GiveHandCard(c.Controller, others[a[1]], hand[a[0]])
			}, engine.HandQuestion(cardkit.AskGiveWhichLootCard), engine.Question{Text: fetalHauntGiveItTo, Options: func(c *engine.Ctx, a []int) []engine.Msg {
				if a[0] < 0 {
					return nil
				}
				return playerLabels(c.G, otherPlayers(c))
			}})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When you die, before paying penalties, give this to another player.",
			Trigger: engine.WhenYouDie(),
			Effects: []engine.Effect{cardkit.GiveThisAway},
		},
	},
}
