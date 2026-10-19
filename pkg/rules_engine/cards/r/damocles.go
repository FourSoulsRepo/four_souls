package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Damocles (Passive Treasure Card)
//
//	If you would gain any number of treasures, instead gain that many +1.
//	When you die, before paying penalties, destroy this.
var damocles = engine.CardDef{
	Ref:    "damocles",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Replacements: []engine.Replacement{{
		Text: "If you would gain any number of treasures, instead gain that many +1.",
		When: func(g *engine.Game, self engine.ObjectID, a engine.Action) bool {
			return a.Kind == engine.ActGainTreasure && a.Player == g.Object(self).Controller
		},
		Do: func(_ *engine.Game, _ engine.ObjectID, a engine.Action) []engine.Action {
			a.Amount++
			return []engine.Action{a}
		},
	}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When you die, before paying penalties, destroy this.",
		Trigger: engine.WhenYouDie(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			c.G.DestroyObject(c.G.Object(c.Source).Controller, c.Source)
		})},
	}},
}
