package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Haunted! (Room Card)
//
//	Each time a player gains a soul, they skip their next turn.
var haunted = engine.CardDef{
	Ref:    "haunted",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player gains a soul, they skip their next turn.",
		Trigger: engine.Trigger{On: engine.EvGainedSoul, Match: anyPlayer},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.Players[c.EventPlayer].SkipTurns++ })},
	}},
}
