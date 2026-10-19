package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// X-Ray Vision (Active Treasure Card)
//
//	{Tap Effect}Look at a player's hand and the top card of a deck. You may put the top card of a deck into discard.
//	At the end of your turn, recharge this.
var xRayVision = engine.CardDef{
	Ref:    "x_ray_vision",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Look at a player's hand and the top card of a deck. You may put the top card of a deck into discard.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.LookAt(c.Controller, c.G.Players[c.Targets[0].Player].Hand...)
			}), cardkit.LookMayDiscard},
		},
		cardkit.RechargeAtEndOfTurn,
	},
}
