package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pieces many Base Game cards share.

// extraLoot is every Base Game character's ↷ ability. It can be used on
// any player's turn, in response to anything.
var extraLoot = engine.Ability{
	Kind:    engine.Activated,
	Text:    "↷: Play an additional loot card this turn.",
	Costs:   []engine.Cost{engine.Tap()},
	Effects: []engine.Effect{engine.AddLootPlays(1)},
}

// rechargeAtEndOfTurn: "At the end of your turn, recharge this".
var rechargeAtEndOfTurn = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "At the end of your turn, recharge this.",
	Trigger: engine.AtEndOfYourTurn(),
	Effects: []engine.Effect{engine.RechargeSelf()},
}

// discardOne is "Discard 1 loot card." for the controller.
var discardOne = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.DiscardFromHand(c.Controller, c.HandCard(a[0]))
	}
}, engine.HandQuestion("Discard which loot card?"))

// damageEqualToRoll deals damage equal to the roll to target 0.
func damageEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.DealDamage(r, 0))
	}
	return t
}
