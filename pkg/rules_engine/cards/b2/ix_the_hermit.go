package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// IX. The Hermit (Wildcard Card)
//
//	Look at the top 5 cards of the treasure deck. Put 1 on top and the rest on the bottom.
var ixTheHermit = engine.CardDef{
	Ref:    "ix_the_hermit",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Look at the top 5 cards of the treasure deck. Put 1 on top and the rest on the bottom.",
			Effects: []engine.Effect{oneOnTop(engine.TreasureDeck)},
		},
	},
}
