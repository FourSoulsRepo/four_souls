package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Pieces many cards share; they live in cardkit for every set.
var (
	itemsOf        = cardkit.ItemsOf
	freeSlots      = cardkit.FreeSlots
	notAttacked    = cardkit.NotAttacked
	inShop         = cardkit.InShop
	putBackInOrder = cardkit.PutBackInOrder
	otherPlayers   = cardkit.OtherPlayers
	damageAMonster = cardkit.DamageAMonster
	playerLabel    = cardkit.PlayerLabel
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

// playerLabels names players as the prompts do.
func playerLabels(g *engine.Game, ps []engine.PlayerID) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = playerLabel(g, p)
	}
	return out
}

// othersInOrder lists the other players in turn order from the
// controller.
func othersInOrder(c *engine.Ctx) []engine.PlayerID {
	var out []engine.PlayerID
	p := c.G.PlayerAfter(c.Controller)
	for range len(c.G.Players) - 1 {
		out = append(out, p)
		p = c.G.PlayerAfter(p)
	}
	return out
}

// eachOtherMay asks each other player, in turn order, a yes-or-no
// question; do gets the players who said yes ("each other player may
// choose to loot 1").
func eachOtherMay(question string, do func(c *engine.Ctx, yes []engine.PlayerID)) engine.Effect {
	qs := make([]engine.Question, 3) // at most 3 other players
	for i := range qs {
		qs[i] = engine.Question{
			Text:   question,
			Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return othersInOrder(c)[i] },
			Options: func(c *engine.Ctx, _ []int) []string {
				if i >= len(othersInOrder(c)) {
					return nil
				}
				return []string{"yes", "no"}
			},
		}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		var yes []engine.PlayerID
		for i, p := range othersInOrder(c) {
			if i < len(a) && a[i] == 0 {
				yes = append(yes, p)
			}
		}
		do(c, yes)
	}, qs...)
}

// eachGivesYouALoot: each of ps gives the controller a loot card of
// their choice, in order.
func eachGivesYouALoot(ps []engine.PlayerID) engine.Effect {
	qs := make([]engine.Question, len(ps))
	for i, p := range ps {
		qs[i] = engine.Question{
			Text:    "Give which loot card?",
			Player:  func(*engine.Ctx, []int) engine.PlayerID { return p },
			Options: func(c *engine.Ctx, _ []int) []string { return labelsOf(c.G, c.G.Players[p].Hand) },
		}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		for i, p := range ps {
			if hand := c.G.Players[p].Hand; a[i] >= 0 && a[i] < len(hand) {
				c.G.GiveHandCard(p, c.Controller, hand[a[i]])
			}
		}
	}, qs...)
}

// discardAShopItem puts a shop item of the controller's choice into
// discard; its slot is refilled.
var discardAShopItem = shopItemToDiscard(false)

// mayDiscardAShopItem is discardAShopItem with "none".
var mayDiscardAShopItem = shopItemToDiscard(true)

func shopItemToDiscard(may bool) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if items := c.G.ShopItems(); a[0] >= 0 && a[0] < len(items) {
			c.G.DiscardShopItem(items[a[0]])
		}
	}, engine.Question{Text: "Put which shop item into discard?", Options: func(c *engine.Ctx, _ []int) []string {
		out := labelsOf(c.G, c.G.ShopItems())
		if may && len(out) > 0 {
			out = append(out, "none")
		}
		return out
	}})
}

// cardsInPlay lists the non-eternal cards in play: characters, items,
// shop items, monsters and rooms.
func cardsInPlay(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, pl := range g.Players {
		out = append(out, pl.Character)
		for _, id := range pl.InPlay {
			if g.Object(id).Role == engine.RoleItem {
				out = append(out, id)
			}
		}
	}
	out = append(out, g.ShopItems()...)
	for _, rows := range [][]engine.Slot{g.Monsters, g.Rooms} {
		for _, s := range rows {
			if top, ok := s.TopOf(); ok {
				out = append(out, top)
			}
		}
	}
	kept := out[:0]
	for _, id := range out {
		if !g.Eternal(id) {
			kept = append(kept, id)
		}
	}
	return kept
}

// changeACounter: "Choose a non-eternal card in play. Put a counter on it
// of a type already on it or remove a counter from it" (Marbles). With
// items, only items, and a generic counter may always be added
// (Playdough Cookie); with none, the controller may choose nothing.
func changeACounter(items bool) engine.Effect {
	cards := func(g *engine.Game) []engine.ObjectID {
		var out []engine.ObjectID
		for _, id := range cardsInPlay(g) {
			o := g.Object(id)
			if items && o.Role != engine.RoleItem {
				continue
			}
			if items || len(o.Counters) > 0 {
				out = append(out, id)
			}
		}
		return out
	}
	type change struct {
		name string
		n    int
	}
	changes := func(g *engine.Game, id engine.ObjectID) []change {
		o := g.Object(id)
		var out []change
		generic := false
		for _, ct := range o.Counters {
			out = append(out, change{ct.Name, 1}, change{ct.Name, -1})
			generic = generic || ct.Name == ""
		}
		if items && !generic {
			out = append(out, change{"", 1})
		}
		return out
	}
	name := func(n string) string {
		if n == "" {
			return "counter"
		}
		return n + " counter"
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		ids := cards(c.G)
		if a[0] < 0 || a[0] >= len(ids) || a[1] < 0 {
			return
		}
		ch := changes(c.G, ids[a[0]])[a[1]]
		c.G.AddNamedCounters(ids[a[0]], ch.name, ch.n)
	}, engine.Question{Text: "Change the counters of which card?", Options: func(c *engine.Ctx, _ []int) []string {
		out := labelsOf(c.G, cards(c.G))
		if items && len(out) > 0 {
			out = append(out, "none")
		}
		return out
	}}, engine.Question{Text: "Which change?", Options: func(c *engine.Ctx, a []int) []string {
		ids := cards(c.G)
		if a[0] < 0 || a[0] >= len(ids) {
			return nil
		}
		var out []string
		for _, ch := range changes(c.G, ids[a[0]]) {
			if ch.n > 0 {
				out = append(out, "add a "+name(ch.name))
			} else {
				out = append(out, "remove a "+name(ch.name))
			}
		}
		return out
	}})
}

// killYou kills the controller ("You die").
var killYou = engine.EffectFunc(func(c *engine.Ctx) {
	if !c.G.Players[c.Controller].Dead {
		them := *c
		them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: c.Controller}}
		them.Do(engine.Kill(0))
	}
})

// discardHandLoot is "discard your hand and loot n".
func discardHandLoot(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, id := range append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...) {
			c.G.DiscardFromHand(c.Controller, id)
		}
		c.Do(engine.Loot(n))
	})
}
