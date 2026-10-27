package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Book Of Virtues (Eternal Treasure Card)
//
//	Each time you deal combat damage, roll-
//	5-6: Deal 1 damage to a monster or player.
//	Each time you kill a monster, choose another player. You and that player each loot 1.
//	-Eternal- This can't be destroyed or put into discard.
var bookOfVirtues = engine.CardDef{
	Ref:     "book_of_virtues",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time you deal combat damage, roll- 5-6: Deal 1 damage to a monster or player.",
			Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Text == "combat" && e.By == int(g.Object(self).Controller)+1
			}},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(5, 6, damageAnything(1)))},
		},
		{
			Kind: engine.Triggered,
			Text: "Each time you kill a monster, choose another player. You and that player each loot 1.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return youKilledMonster(g, g.Object(self).Controller, e)
			}},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				others := otherPlayers(c)
				if a[0] < 0 || a[0] >= len(others) {
					return
				}
				c.Do(engine.Loot(1))
				them := *c
				them.Controller = others[a[0]]
				them.Do(engine.Loot(1))
			}, engine.Question{Text: cardkit.AskChooseAnotherPlayer, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				return playerLabels(c.G, otherPlayers(c))
			}})},
		},
	},
}
