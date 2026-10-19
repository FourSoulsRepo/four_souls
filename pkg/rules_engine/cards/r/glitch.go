package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Glitch (Eternal Treasure Card)
//
//	As your turn starts, choose a shop item. This gains the abilities of that item till end of turn.
//	As your turn ends, choose a shop item. This gains the abilities of that item till the start of your next turn. Recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var glitch = engine.CardDef{
	Ref:     "glitch",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "As your turn starts, choose a shop item. This gains the abilities of that item till end of turn.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{gainShopAbilities(engine.GainEndOfTurn)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "As your turn ends, choose a shop item. This gains the abilities of that item till the start of your next turn. Recharge this.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{gainShopAbilities(engine.GainTillYourTurn), engine.RechargeSelf()},
		},
	},
}

// gainShopAbilities: the controller chooses a shop item; this gains its
// abilities for that long.
func gainShopAbilities(till int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if items := c.G.ShopItems(); a[0] >= 0 && a[0] < len(items) {
			c.G.GainAbilities(c.Source, c.G.CardOf(items[a[0]]), till)
		}
	}, engine.Question{Text: "Gain the abilities of which shop item?", Options: func(c *engine.Ctx, _ []int) []string {
		return labelsOf(c.G, c.G.ShopItems())
	}})
}
