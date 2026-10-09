package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// IV. The Emperor (Wildcard Card)
//
//	Look at the top 5 cards of the monster deck. Put 1 on top and the rest on the bottom.
var ivTheEmperor = engine.CardDef{
	Ref:    "iv_the_emperor",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Look at the top 5 cards of the monster deck. Put 1 on top and the rest on the bottom.",
			Effects: []engine.Effect{oneOnTop(engine.MonsterDeck)},
		},
	},
}
