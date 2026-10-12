package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Credit Card (Wildcard Card)
//
//	The first shop item a player purchases this turn costs 0¢.
var creditCard = engine.CardDef{
	Ref:    "credit_card",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "The first shop item a player purchases this turn costs 0¢.",
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Turn.FreeShopItem = true // counts only if no shop item was bought yet this turn
		})},
	}},
}
