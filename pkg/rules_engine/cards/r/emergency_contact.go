package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Emergency Contact (Wildcard Card)
//
//	Cancel an attack on a monster and put that monster card on the bottom of the monster deck.
var emergencyContact = engine.CardDef{
	Ref:    "emergency_contact",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Cancel an attack on a monster and put that monster card on the bottom of the monster deck.",
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			g := c.G
			// A deck attack counts once a monster is revealed (R-ATK-05, Q-24).
			if !g.Attack.Started || g.Object(g.Attack.Target).Role != engine.RoleMonster {
				return
			}
			m := g.Attack.Target
			c.Do(engine.CancelAttack())
			g.SlotToDeckBottom(m)
			g.RefillSlots()
		})},
	}},
}
