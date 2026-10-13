package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gold Key (Wildcard Card)
//
//	The active player may attack the monster deck any number of times till end of turn.
var goldKey = engine.CardDef{
	Ref:    "gold_key",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "The active player may attack the monster deck any number of times till end of turn.",
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			for range 99 { // "any number of times"
				c.G.AddDeckAttack()
			}
		})},
	}},
}
