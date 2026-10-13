package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Missing Page (Passive Treasure Card)
//
//	Each time a player would roll a ❺, you may force that player to reroll it.
var theMissingPage = engine.CardDef{
	Ref:    "the_missing_page",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player would roll a 5, you may force that player to reroll it.",
		Trigger: engine.WhenARollWouldBe(5),
		Effects: []engine.Effect{may("Force a reroll of the 5?", engine.EffectFunc(func(c *engine.Ctx) { c.G.Reroll(c.EventStack) }))},
	}},
}
