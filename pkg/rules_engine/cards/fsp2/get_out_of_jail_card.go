package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Get Out Of Jail Card (Wildcard Card)
//
//	Other players can't play loot cards or activate items till end of turn.
var getOutOfJailCard = engine.CardDef{
	Ref:    "get_out_of_jail_card",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Other players can't play loot cards or activate items till end of turn.",
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatLockTurn, Player: c.Controller, Amount: 1})
		})},
	}},
}
