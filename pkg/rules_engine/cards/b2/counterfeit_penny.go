package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Counterfeit Penny (Trinket Card)
//
//	If you would gain any number of ¢, gain that much +1¢ instead.
//	-Trinket- This loot becomes an item under your control when it resolves.
var counterfeitPenny = engine.CardDef{
	Ref:     "counterfeit_penny",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Replacements: []engine.Replacement{{
		Text: "If you would gain any number of ¢, gain that much +1¢ instead.",
		When: func(g *engine.Game, self engine.ObjectID, a engine.Action) bool {
			return a.Kind == engine.ActGainCents && a.Amount > 0 && a.Player == g.Object(self).Controller
		},
		Do: func(_ *engine.Game, _ engine.ObjectID, a engine.Action) []engine.Action {
			a.Amount++
			return []engine.Action{a}
		},
	}},
	Abilities: []engine.Ability{},
}
