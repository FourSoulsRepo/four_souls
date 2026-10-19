package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Two Of Spades (Wildcard Card)
//
//	Choose another player. Loot x, where x is the number of loot cards in that player's hand.
var twoOfSpades = engine.CardDef{
	Ref:     "two_of_spades",
	Kind:    engine.LootCard,
	Copies:  1,
	Rewards: nil,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Choose another player. Loot x, where x is the number of loot cards in that player's hand.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if n := len(c.G.Players[c.Targets[0].Player].Hand); n > 0 {
				c.Do(engine.Loot(n))
			}
		})},
	}},
}
