package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Overflow! (Bad Event Card)
//
//	Each player rolls-
//	1-2: They skip their next turn.
var overflow = engine.CardDef{
	Ref:    "overflow",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each player rolls- 1-2: They skip their next turn.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.EachPlayer(engine.Roll(engine.RollTable{}.Results(1, 2, engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Players[c.Controller].SkipTurns++
		}))))},
	}},
}
