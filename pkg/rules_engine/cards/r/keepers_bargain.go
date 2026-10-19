package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper’s Bargain (Eternal Treasure Card)
//
//	{Tap Effect}Each other player may choose to gain 1¢. Gain 1¢ + 1¢ for each player who did.
//	Shop items you purchase cost 3¢ less.
//	-Eternal- This can't be destroyed or put into discard.
var keepersBargain = engine.CardDef{
	Ref:     "keepers_bargain",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Statics: []engine.Static{engine.YouHave(engine.StatShopPrice, -3)},
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Each other player may choose to gain 1¢. Gain 1¢ + 1¢ for each player who did.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{eachOtherMay("Gain 1¢? They gain 1¢ more.", func(c *engine.Ctx, yes []engine.PlayerID) {
			for _, p := range yes {
				them := *c
				them.Controller = p
				them.Do(engine.GainCents(1))
			}
			c.Do(engine.GainCents(1 + len(yes)))
		})},
	}},
}
