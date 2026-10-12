package g2

import (
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Pieces many cards share; they live in cardkit for every set.
var (
	extraLoot           = cardkit.ExtraLoot
	rechargeAtEndOfTurn = cardkit.RechargeAtEndOfTurn
	livingPlayers       = cardkit.LivingPlayers
	rechargeAnItem      = cardkit.RechargeAnItem
)

// centsEqualToRoll: "that player gains ¢ equal to the result" for
// target 0.
func centsEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.EffectFunc(func(c *engine.Ctx) {
			them := *c
			them.Controller = c.Targets[0].Player
			them.Do(engine.GainCents(r))
		}))
	}
	return t
}

// notAttacked: a monster that is not the one being attacked.
func notAttacked(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	return !g.Attack.On || g.Attack.Target != c.Object
}

// inShop: an item in a shop slot.
func inShop(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	return slices.Contains(g.ShopItems(), c.Object)
}

// discardTarget puts target 0, a monster or a shop item, into discard;
// its slot is refilled.
var discardTarget = engine.EffectFunc(func(c *engine.Ctx) {
	id := c.Targets[0].Object
	if slices.Contains(c.G.ShopItems(), id) {
		c.G.DiscardShopItem(id) // refills its slot
		return
	}
	c.G.DiscardMonster(id)
	c.G.RefillSlots()
})

// pickAnotherPlayer asks the controller for another living player.
func pickAnotherPlayer(text string) engine.Question {
	return engine.Question{Text: text, Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, p := range livingPlayers(c.G, c.Controller) {
			out = append(out, cardkit.PlayerLabel(c.G, p))
		}
		return out
	}}
}

// forTarget runs effects as if target 0, a player, controlled them:
// "they gain …".
func forTarget(effects ...engine.Effect) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		them := *c
		them.Controller = c.Targets[0].Player
		them.Do(effects...)
	})
}

// rechargeYourItems: "Recharge each item you control".
var rechargeYourItems = engine.EffectFunc(func(c *engine.Ctx) {
	for _, id := range c.G.Players[c.Controller].InPlay {
		if c.G.Object(id).Role == engine.RoleItem {
			c.G.Recharge(id)
		}
	}
})

// attackedDC: "Monsters you attack have +n DC till end of turn".
func attackedDC(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatAttackedDC, Player: c.Controller, Amount: n})
	})
}

// combatDamageOnRollOf triggers when the controller takes combat damage
// from a monster on an attack roll of n; c.EventSource is the monster.
func combatDamageOnRollOf(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "combat" && e.Player == g.Object(self).Controller && g.Attack.LastRoll == n
	}}
}
