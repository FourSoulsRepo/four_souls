package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ace Of Diamonds (Wildcard Card)
//
//	Gain x¢, where x is the number of monster slots plus the number of loot cards in your hand.
var aceOfDiamonds = engine.CardDef{
	Ref:    "ace_of_diamonds",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Gain x¢, where x is the number of monster slots plus the number of loot cards in your hand.",
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.Do(engine.GainCents(len(c.G.Monsters) + len(c.G.Players[c.Controller].Hand)))
		})},
	}},
}
