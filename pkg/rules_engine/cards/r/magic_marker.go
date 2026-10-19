package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Magic Marker (Wildcard Card)
//
//	Change a number in the effect text of a card in play or loot being played by 1 till end of turn. The number can't go below 1 or above 6.
var magicMarker = engine.CardDef{
	Ref:    "magic_marker",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Change a number in the effect text of a card in play or loot being played by 1 till end of turn. The number can't go below 1 or above 6.",
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			ids := markable(c)
			if a[0] < 0 || a[0] >= len(ids) {
				return
			}
			ns := c.G.NumbersOf(ids[a[0]])
			if a[1] < 0 || a[1] >= len(ns) || a[2] < 0 {
				return
			}
			delta := 1
			if a[2] == 1 {
				delta = -1
			}
			c.G.MarkNumber(ids[a[0]], ns[a[1]], delta)
		}, engine.Question{Text: "Change a number of which card?", Options: func(c *engine.Ctx, _ []int) []string {
			return labelsOf(c.G, markable(c))
		}}, engine.Question{Text: "Which number?", Options: func(c *engine.Ctx, a []int) []string {
			ids := markable(c)
			if a[0] < 0 || a[0] >= len(ids) {
				return nil
			}
			var out []string
			for _, n := range c.G.NumbersOf(ids[a[0]]) {
				out = append(out, n.Label())
			}
			return out
		}}, engine.Question{Text: "Change it by?", Options: func(_ *engine.Ctx, a []int) []string {
			if a[1] < 0 {
				return nil
			}
			return []string{"+1", "-1"}
		}})},
	}},
}

// markable lists the cards in play and loot cards being played that have
// numbers Magic Marker can change (R-ABIL-39).
func markable(c *engine.Ctx) []engine.ObjectID {
	var all []engine.ObjectID
	for _, pl := range c.G.Players {
		all = append(all, pl.Character)
		all = append(all, pl.InPlay...)
	}
	for _, rows := range [][]engine.Slot{c.G.Monsters, c.G.Rooms} {
		for _, s := range rows {
			if top, ok := s.TopOf(); ok {
				all = append(all, top)
			}
		}
	}
	for _, it := range c.G.Stack {
		if it.Kind == engine.StackLoot && it.Source != c.Source {
			all = append(all, it.Source)
		}
	}
	var out []engine.ObjectID
	for _, id := range all {
		if len(c.G.NumbersOf(id)) > 0 {
			out = append(out, id)
		}
	}
	return out
}
