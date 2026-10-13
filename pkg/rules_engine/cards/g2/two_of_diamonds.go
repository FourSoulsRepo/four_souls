package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Two Of Diamonds (Wildcard Card)
//
//	Choose a player. They gain x¢, where x is the number of ¢ they have.
var twoOfDiamonds = engine.CardDef{
	Ref:    "two_of_diamonds",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Choose a player. They gain x¢, where x is the number of ¢ they have.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
		Effects: []engine.Effect{forTarget(engine.EffectFunc(func(c *engine.Ctx) {
			c.Do(engine.GainCents(c.G.Players[c.Controller].Cents))
		}))},
	}},
}
