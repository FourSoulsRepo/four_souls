package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Impulse (Curse Card)
//
//	{Curse Effect}At the end of your turn, deactivate each item you control and your character.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfImpulse = engine.CardDef{
	Ref:    "curse_of_impulse",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of your turn, deactivate each item you control and your character.",
		Trigger: engine.AtEndOfYourTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Object(c.G.Players[c.Controller].Character).Charged = false
			for _, id := range itemsOf(c.G, c.Controller) {
				c.G.Object(id).Charged = false
			}
		})},
	}},
}
