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
