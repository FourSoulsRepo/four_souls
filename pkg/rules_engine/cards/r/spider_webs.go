package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spider Webs (Room Card)
//
//	At the end of the turn, the active player deactivates their character.
var spiderWebs = engine.CardDef{
	Ref:    "spider_webs",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of the turn, the active player deactivates their character.",
		Trigger: atEndOfTheTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.Object(c.G.Players[c.Controller].Character).Charged = false
		})},
	}},
}
