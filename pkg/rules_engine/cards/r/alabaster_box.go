package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Alabaster Box (Passive Treasure Card)
//
//	At the start of your turn, put a counter on this.
//	When the 2nd counter is put on this, gain 10¢.
//	When the 4th counter is put on this, loot 4.
//	When the 6th counter is put on this, gain +2 treasure.
var alabasterBox = engine.CardDef{
	Ref:    "alabaster_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, put a counter on this.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{Kind: engine.Triggered, Text: "When the 2nd counter is put on this, gain 10¢.", Trigger: whenCounterReaches(2), Effects: []engine.Effect{engine.GainCents(10)}},
		{Kind: engine.Triggered, Text: "When the 4th counter is put on this, loot 4.", Trigger: whenCounterReaches(4), Effects: []engine.Effect{engine.Loot(4)}},
		{Kind: engine.Triggered, Text: "When the 6th counter is put on this, gain +2 treasure.", Trigger: whenCounterReaches(6), Effects: []engine.Effect{engine.GainTreasure(2)}},
	},
}
