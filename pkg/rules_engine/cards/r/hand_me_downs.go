package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	handMeDownsLoot1If = engine.NewKey("card.hand_me_downs.loot_1_if_you_do")
)

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
		Effects: []engine.Effect{eachOtherMay(handMeDownsLoot1If, func(c *engine.Ctx, yes []engine.PlayerID) {
			for _, p := range yes {
				them := *c
				them.Controller = p
				them.Do(engine.Loot(1))
			}
			c.Do(eachGivesYouALoot(yes))
		})},
	}},
}
