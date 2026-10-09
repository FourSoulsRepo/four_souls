package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Map (Passive Treasure Card)
//
//	At the end of your turn, look at the top 4 cards of the monster deck. Put them back in any order.
var theMap = engine.CardDef{
	Ref:    "the_map",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, look at the top 4 cards of the deck. Put them back in any order.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{putBackInOrder(engine.MonsterDeck, 4)},
		},
	},
}
