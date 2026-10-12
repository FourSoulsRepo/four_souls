package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Jera (Pill/Rune Card)
//
//	Loot x, where x is equal to the number of loot cards in your hand.
var jera = engine.CardDef{
	Ref:    "jera",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Loot x, where x is equal to the number of loot cards in your hand.",
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.Do(engine.Loot(len(c.G.Players[c.Controller].Hand)))
		})},
	}},
}
