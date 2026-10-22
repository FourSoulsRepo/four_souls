package r

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Betrayal! (Good Event Card)
//
//	Choose another player. Roll-
//	1-2: Steal 5¢ from them.
//	3-4: Steal a loot card from them at random.
//	5: Steal a non-eternal item they control.
//	6: Steal a soul they control.
var betrayal = engine.CardDef{
	Ref:    "betrayal",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Choose another player. Roll- 1-2: Steal 5¢ from them. 3-4: Steal a loot card from them at random. 5: Steal a non-eternal item they control. 6: Steal a soul they control.",
		Takes:   engine.TakesOthers,
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.AskWithRoll(func(c *engine.Ctx, a []int) {
			if ps := othersOfActive(c.G); a[0] >= 0 && a[0] < len(ps) {
				them := *c
				them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: ps[a[0]]}}
				them.Do(engine.Roll(betrayalTable))
			}
		}, betrayalTable, engine.Question{Text: "Choose another player.", Options: func(c *engine.Ctx, _ []int) []string {
			return playerLabels(c.G, othersOfActive(c.G))
		}})},
	}},
}

var betrayalTable = engine.RollTable{}.
	Results(1, 2, engine.StealCents(5, 0)).
	Results(3, 4, engine.Ask(func(c *engine.Ctx, a []int) {
		from := c.Targets[0].Player
		if hand := c.G.Players[from].Hand; a[0] >= 0 && a[0] < len(hand) {
			c.G.GiveHandCard(from, c.Controller, hand[a[0]])
		}
	}, engine.Question{Text: "A loot card at random", Random: true, Options: func(c *engine.Ctx, _ []int) []string {
		// Not the cards' names: the pick is public, the hand is not (A-16).
		var out []string
		for i := range c.G.Players[c.Targets[0].Player].Hand {
			out = append(out, "card "+strconv.Itoa(i+1))
		}
		return out
	}})).
	Results(5, 5, stealOf(func(g *engine.Game, id engine.ObjectID) bool {
		return g.Object(id).Role == engine.RoleItem && !g.Eternal(id)
	})).
	Results(6, 6, stealOf(func(g *engine.Game, id engine.ObjectID) bool { return g.Object(id).Role == engine.RoleSoul }))

// stealOf: the controller steals an object of target 0 that keep picks.
func stealOf(keep func(g *engine.Game, id engine.ObjectID) bool) engine.Effect {
	objects := func(c *engine.Ctx) []engine.ObjectID {
		var out []engine.ObjectID
		for _, id := range c.G.Players[c.Targets[0].Player].InPlay {
			if keep(c.G, id) {
				out = append(out, id)
			}
		}
		return out
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if ids := objects(c); a[0] >= 0 && a[0] < len(ids) {
			c.G.GainControl(c.Controller, ids[a[0]])
		}
	}, engine.Question{Text: "Steal which?", Options: func(c *engine.Ctx, _ []int) []string { return labelsOf(c.G, objects(c)) }})
}
