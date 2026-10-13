package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Eye (Passive Treasure Card)
//
//	Combat damage you deal on attack rolls of 6 is increased by 3.
//	When you roll an attack roll of 1, end your turn. Cancel everything that hasn't resolved.
var cursedEye = engine.CardDef{
	Ref:    "cursed_eye",
	Kind:   engine.TreasureCard,
	Copies: 1,
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if hit && roll == 6 {
			return n + 3 // "Combat damage you deal on attack rolls of 6 is increased by 3."
		}
		return n
	},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When you roll an attack roll of 1, end your turn. Cancel everything that hasn't resolved.",
		Trigger: onYourAttackRollOf(1),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.EndTurnNow() })},
	}},
}
