package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ehwaz (Pill/Rune Card)
//
//	Put each monster not being attacked into discard and replace each with the top card of the monster deck.
var ehwaz = engine.CardDef{
	Ref:    "ehwaz",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind: engine.LootAbility,
			Text: "Put each monster not being attacked into discard and replace each with the top card of the monster deck.",
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				for _, s := range c.G.Monsters {
					if top, ok := s.TopOf(); ok && c.G.Object(top).Role == engine.RoleMonster && top != c.G.Attack.Target {
						c.G.DiscardMonster(top)
					}
				}
				c.G.RefillSlots()
			})},
		},
	},
}
