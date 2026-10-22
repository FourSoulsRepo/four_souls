package fsp2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// P.H.D. (Passive Treasure Card)
//
//	You may add or subtract 1 from any of your non-attack rolls.
var phd = engine.CardDef{
	Ref:    "phd",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "You may add or subtract 1 from any of your non-attack rolls.",
		Trigger: engine.Trigger{On: engine.EvRollWouldResolve, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			it, ok := g.StackItemByID(e.StackID)
			return ok && e.Player == g.Object(self).Controller && !it.Attack && !cardkit.UsedOn(g, self, e.StackID) // once per roll (R-DICE-10)
		}},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			cardkit.MarkUsedOn(c.G, c.Source, c.EventStack)
			it, ok := c.G.StackItemByID(c.EventStack)
			if !ok || a[0] < 0 || a[0] == 2 {
				return
			}
			c.G.SetRoll(it.ID, min(max(it.Roll+[]int{1, -1}[a[0]], 1), 6), c.Source)
		}, engine.Question{Text: "Change your roll?", Options: func(*engine.Ctx, []int) []string { return []string{"+1", "-1", "No"} }})},
	}},
}
