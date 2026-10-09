package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Compass (Passive Treasure Card)
//
//	At the end of your turn, look at the top 4 cards of the loot deck. Put them back in any order.
var theCompass = engine.CardDef{
	Ref:    "the_compass",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, look at the top 4 cards of the deck. Put them back in any order.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{putBackInOrder(engine.LootDeck, 4)},
		},
	},
}
