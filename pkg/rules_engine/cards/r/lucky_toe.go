package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lucky Toe (Trinket Card)
//
//	If you would loot, except during the loot step, instead loot that much +1.
//	-Trinket- This loot becomes an item under your control when it resolves.
var luckyToe = engine.CardDef{
	Ref:     "lucky_toe",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Replacements: []engine.Replacement{{
		Text: "If you would loot, except during the loot step, instead loot that much +1.",
		When: func(g *engine.Game, self engine.ObjectID, a engine.Action) bool {
			return a.Kind == engine.ActLoot && a.Player == g.Object(self).Controller &&
				(g.Turn.Step != engine.StepLoot || a.Player != g.Turn.Active)
		},
		Do: func(_ *engine.Game, _ engine.ObjectID, a engine.Action) []engine.Action {
			a.Amount++
			return []engine.Action{a}
		},
	}},
}
