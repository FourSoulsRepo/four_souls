package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spindown Dice (Eternal Treasure Card)
//
//	{Tap Effect}Destroy an item. If it was controlled by a player, they may steal an item from the shop.
//	At the end of your turn, roll-
//	4-6: Recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var spindownDice = engine.CardDef{
	Ref:       "spindown_dice",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
	Tap:       true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Destroy a non-eternal item. If it was controlled by a player, they may steal an item from the shop.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, func(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool { return !g.Eternal(c.Object) })},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				id := c.Targets[0].Object
				p := c.G.Object(id).Controller
				// The steal needs the destruction (Q-16, owner).
				if c.G.DestroyObject(p, id) && p != engine.NoPlayer {
					c.Do(stealFromShop(p, true))
				}
			})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, roll- 4-6: Recharge this.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(4, 6, engine.RechargeSelf()))},
		},
	},
}
