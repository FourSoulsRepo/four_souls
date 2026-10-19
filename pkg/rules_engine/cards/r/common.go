package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Pieces many cards share; they live in cardkit for every set.
var (
	itemsOf   = cardkit.ItemsOf
	freeSlots = cardkit.FreeSlots
)

// eachOtherThan deals n damage from the controller to each player but
// skip and the controller, in turn order from the controller (R-MECH-28).
func eachOtherThan(c *engine.Ctx, skip engine.PlayerID, n int) {
	order := []engine.PlayerID{}
	p := c.Controller
	for range c.G.Players {
		if p != skip && p != c.Controller {
			order = append(order, p)
		}
		p = c.G.PlayerAfter(p)
	}
	// The first in turn order resolves first: it goes on the stack last.
	for i := len(order) - 1; i >= 0; i-- {
		c.G.DealDamageTo(engine.Target{Player: order[i], IsPlayer: true}, n, c.Controller, c.Source)
	}
}

// stealFromShop: player p steals a shop item of their choice; with may
// they may decline.
func stealFromShop(p engine.PlayerID, may bool) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if items := c.G.ShopItems(); a[0] >= 0 && a[0] < len(items) {
			c.G.GainControl(p, items[a[0]])
		}
	}, engine.Question{
		Text:   "Steal which shop item?",
		Player: func(*engine.Ctx, []int) engine.PlayerID { return p },
		Options: func(c *engine.Ctx, _ []int) []string {
			out := labelsOf(c.G, c.G.ShopItems())
			if may && len(out) > 0 {
				out = append(out, "none")
			}
			return out
		},
	})
}

// labelsOf lists the cards of objects.
func labelsOf(g *engine.Game, ids []engine.ObjectID) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = string(g.Object(id).Card)
	}
	return out
}
