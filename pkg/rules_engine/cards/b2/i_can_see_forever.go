package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// I Can See Forever! (Good Event Card)
//
//	Look at the top 6 cards of the loot deck. Put them back in any order, then loot 1.
var iCanSeeForever = engine.CardDef{
	Ref:    "i_can_see_forever",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Look at the top 6 cards of the loot deck. Put them back in any order, then loot 1.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{putBackInOrder(engine.LootDeck, 6), engine.Loot(1)},
		},
	},
}
