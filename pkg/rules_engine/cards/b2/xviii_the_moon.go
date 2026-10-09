package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XVIII. The Moon (Wildcard Card)
//
//	Look at the top 5 cards of the loot deck. Put 1 on top and the rest on the bottom.
var xviiiTheMoon = engine.CardDef{
	Ref:    "xviii_the_moon",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Look at the top 5 cards of the loot deck. Put 1 on top and the rest on the bottom.",
			Effects: []engine.Effect{oneOnTop(engine.LootDeck)},
		},
	},
}
