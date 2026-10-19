package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Undefined (Active Treasure Card)
//
//	{Tap Effect}Choose a shop item. This becomes a copy of that item till the start of your next turn. Recharge this.
//	At the start of your turn you may put a shop item into discard.
var undefined = engine.CardDef{
	Ref:    "undefined",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a shop item. This becomes a copy of that item till the start of your next turn. Recharge this.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, inShop)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				o := c.G.Object(c.Source)
				if o.Zone.Kind != engine.ZoneInPlay {
					return
				}
				o.CopyOf, o.CopyThisTurn, o.CopyTillYourTurn = c.G.CardOf(c.Targets[0].Object), false, true
			}), engine.RechargeSelf()},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn you may put a shop item into discard.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{mayDiscardAShopItem},
		},
	},
}
