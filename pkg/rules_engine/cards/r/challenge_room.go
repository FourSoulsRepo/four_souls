package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Challenge Room (Room Card)
//
//	Each time a monster dies, the active player rolls-
//	4-6: They must make an additional attack this turn.
var challengeRoom = engine.CardDef{
	Ref:    "challenge_room",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a monster dies, the active player rolls- 4-6: They must make an additional attack this turn.",
		Trigger: engine.WhenAMonsterDies(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(4, 6, engine.EffectFunc(func(c *engine.Ctx) {
			c.G.ForceAttacks(1, false)
		})))},
	}},
}
