package r

import engine "github.com/FourSoulsRepo/rules_engine"

// A Lucky Penny! (1 Cent Card)
//
//	Gain 1¢.
//	Rewards are doubled till end of turn.
var aLuckyPenny = engine.CardDef{
	Ref:    "a_lucky_penny",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Gain 1¢. Rewards are doubled till end of turn.",
		Effects: []engine.Effect{engine.GainCents(1), engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatDoubleRewards, Player: engine.NoPlayer, Amount: 1})
		})},
	}},
}
