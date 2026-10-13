package fsp2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Pieces many cards share; they live in cardkit for every set.
var (
	extraLoot      = cardkit.ExtraLoot
	discardOne     = cardkit.DiscardOne
	rechargeAnItem = cardkit.RechargeAnItem
	damageAPlayer  = cardkit.DamageAPlayer
	damageAMonster = cardkit.DamageAMonster
)

// whenAnotherPlayerDies triggers each time another player dies.
func whenAnotherPlayerDies() engine.Trigger {
	return engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Player != engine.NoPlayer && e.Player != g.Object(self).Controller
	}}
}

// reorderADeck is "Look at the top n cards of a deck and put them back in
// any order": the deck first, then the cards from the top.
func reorderADeck(n int) engine.Effect {
	left := func(c *engine.Ctx, a []int) []engine.ObjectID {
		cards := c.G.DeckTop(c.Deck(a[0]), n)
		for _, i := range a[1:] {
			if i < 0 || i >= len(cards) {
				break
			}
			cards = append(cards[:i:i], cards[i+1:]...)
		}
		return cards
	}
	qs := []engine.Question{engine.DeckQuestion("Look at the top cards of which deck?")}
	for range n {
		qs = append(qs, engine.Question{
			Text: "Which card goes next, from the top?",
			Options: func(c *engine.Ctx, a []int) []string {
				if a[0] < 0 {
					return nil
				}
				var out []string
				for _, id := range left(c, a) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			},
		})
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] < 0 {
			return
		}
		d := c.Deck(a[0])
		cards := c.G.DeckTop(d, n)
		c.G.LookAt(c.Controller, cards...)
		var order []engine.ObjectID
		for _, i := range a[1:] {
			if i < 0 || i >= len(cards) {
				break
			}
			order = append(order, cards[i])
			cards = append(cards[:i:i], cards[i+1:]...)
		}
		c.G.SetDeckTop(d, append(order, cards...))
	}, qs...)
}

// itemsOf lists the items p controls, or every player's for NoPlayer.
func itemsOf(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, pl := range g.Players {
		if p != engine.NoPlayer && pl.ID != p {
			continue
		}
		for _, id := range pl.InPlay {
			if g.Object(id).Role == engine.RoleItem {
				out = append(out, id)
			}
		}
	}
	return out
}

// rerollAnItem is "Reroll an item" (yours only, if set), picked on
// resolution; eternal items can't be rerolled (R-MECH-48).
func rerollAnItem(yours bool) engine.Effect {
	items := func(c *engine.Ctx) []engine.ObjectID {
		p := engine.NoPlayer
		if yours {
			p = c.Controller
		}
		var out []engine.ObjectID
		for _, id := range itemsOf(c.G, p) {
			if !c.G.Eternal(id) {
				out = append(out, id)
			}
		}
		return out
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if ids := items(c); a[0] >= 0 && a[0] < len(ids) {
			c.G.RerollItem(ids[a[0]])
		}
	}, engine.Question{Text: "Reroll which item?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range items(c) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	}})
}

// rerollYourItems is "Reroll each item you control".
var rerollYourItems = engine.EffectFunc(func(c *engine.Ctx) {
	for _, id := range itemsOf(c.G, c.Controller) {
		if !c.G.Eternal(id) {
			c.G.RerollItem(id)
		}
	}
})

// discardHandThenLoot is "Discard your hand, then loot n".
func discardHandThenLoot(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, id := range append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...) {
			c.G.DiscardFromHand(c.Controller, id)
		}
		c.Do(engine.Loot(n))
	})
}

// cancelALoot is "Cancel the effect of a loot being played", picked on
// resolution among the loot cards on the stack.
var cancelALoot = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := lootsOnStack(c); a[0] >= 0 && a[0] < len(ids) {
		c.G.CancelStackItem(ids[a[0]])
	}
}, engine.Question{Text: "Cancel which loot card?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range lootsOnStack(c) {
		it, _ := c.G.StackItemByID(id)
		out = append(out, string(it.Card))
	}
	return out
}})

func lootsOnStack(c *engine.Ctx) []int {
	var out []int
	for _, it := range c.G.Stack {
		if it.Kind == engine.StackLoot && it.Source != c.Source {
			out = append(out, it.ID)
		}
	}
	return out
}

// whenYouMiss triggers when the controller's attack roll misses.
func whenYouMiss() engine.Trigger {
	return engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "attack" && e.Player == g.Object(self).Controller && g.Attack.Started && e.Amount < g.Evasion(g.Attack.Target)
	}}
}
