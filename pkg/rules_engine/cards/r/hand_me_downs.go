package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hand Me Downs (Active Treasure Card)
//
//	{Tap Effect}Each other player may choose to loot 1. Each player that does gives you a loot card.
var handMeDowns = engine.CardDef{
	Ref:    "hand_me_downs",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Each other player may choose to loot 1. Each player that does gives you a loot card.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{eachOtherMay("Loot 1? If you do, you give them a loot card.", func(c *engine.Ctx, yes []engine.PlayerID) {
			for _, p := range yes {
				them := *c
				them.Controller = p
				them.Do(engine.Loot(1))
			}
			c.Do(eachGivesYouALoot(yes))
		})},
	}},
}
