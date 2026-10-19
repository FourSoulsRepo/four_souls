package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blind Rage (Room Card)
//
//	The active player must attack the monster deck once each turn if able.
var blindRage = engine.CardDef{
	Ref:    "blind_rage",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "The active player must attack the monster deck once each turn if able.",
		Trigger: atStartOfTheTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Turn.MustAttacks++
			c.G.Turn.MustAttackDeck++ // the first attack (S-DIS: Blind Rage)
		})},
	}},
}
